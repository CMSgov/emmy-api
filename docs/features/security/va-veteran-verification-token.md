# Getting a VA Access Token with a Client ID and PEM Key

This guide documents the VA client-credentials flow used by the Veteran Service
History and Eligibility APIs when VA issues you:

- a `client_id`
- an RSA private key in PEM format

Unlike a basic `client_id` + `client_secret` exchange, this flow requires you
to build and sign a short-lived JWT, then send that JWT as a
`client_assertion` in the token request.

## What you need

Inputs you already have:

- `client_id`
- a private RSA key PEM file

Values you still need from the VA API docs or your onboarding materials:

- `token_endpoint`: the URL you `POST` to when requesting the token
- `aud`: the JWT audience claim
- `scope`: the scope or scopes allowed for your client

Important: VA's public Postman collection treats `aud` and `token_endpoint` as
separate values. Do not assume they are identical unless the API's
client-credentials documentation says they are.

## JWT contents

VA's public Postman collection signs a JWT with this header:

```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

and this payload shape:

```json
{
  "aud": "<AUD>",
  "iss": "<CLIENT_ID>",
  "sub": "<CLIENT_ID>",
  "jti": "<RANDOM_UUID>",
  "iat": 1712262856,
  "exp": 1712262916
}
```

Notes:

- `iss` and `sub` are both your VA-issued `client_id`
- `jti` is a fresh UUID for each request
- `iat` is the current Unix timestamp in seconds
- `exp` is `iat + 60`
- the JWT is signed with your PEM private key using `RS256`

## Token request

After you sign the JWT, send it to the VA token endpoint as a form-encoded
request:

```text
POST <TOKEN_ENDPOINT>
Content-Type: application/x-www-form-urlencoded
```

Form fields:

- `grant_type=client_credentials`
- `client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer`
- `client_assertion=<SIGNED_JWT>`
- `scope=<SPACE_SEPARATED_SCOPES>`
- `launch=<BASE64_VALUE>` only if the VA API docs say `launch` is required

## Working `curl` example

This example uses `openssl`, `jq`, and `uuidgen` to generate the signed JWT and
exchange it for an access token.

```bash
CLIENT_ID='your-client-id'
PRIVATE_KEY_PEM='/absolute/path/to/private-key.pem'
TOKEN_ENDPOINT='https://sandbox-api.va.gov/oauth2/veteran-verification/system/v1/token'
AUD='https://deptva-eval.okta.com/oauth2/your-auth-server/v1/token'
SCOPE='disability_rating.read veteran_status.read'

NOW="$(date +%s)"
JTI="$(uuidgen | tr '[:upper:]' '[:lower:]')"

b64url() {
  openssl base64 -A | tr '+/' '-_' | tr -d '='
}

HEADER="$(
  printf '%s' '{"alg":"RS256","typ":"JWT"}' | b64url
)"

PAYLOAD="$(
  jq -nc \
    --arg aud "$AUD" \
    --arg iss "$CLIENT_ID" \
    --arg sub "$CLIENT_ID" \
    --arg jti "$JTI" \
    --argjson iat "$NOW" \
    --argjson exp "$((NOW + 60))" \
    '{aud:$aud,iss:$iss,sub:$sub,jti:$jti,iat:$iat,exp:$exp}' | b64url
)"

SIGNING_INPUT="${HEADER}.${PAYLOAD}"

SIGNATURE="$(
  printf '%s' "$SIGNING_INPUT" \
    | openssl dgst -sha256 -sign "$PRIVATE_KEY_PEM" \
    | b64url
)"

CLIENT_ASSERTION="${SIGNING_INPUT}.${SIGNATURE}"

curl --request POST "$TOKEN_ENDPOINT" \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'grant_type=client_credentials' \
  --data-urlencode 'client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer' \
  --data-urlencode "client_assertion=${CLIENT_ASSERTION}" \
  --data-urlencode "scope=${SCOPE}"
```

If your VA API requires a `launch` field, add:

```bash
--data-urlencode 'launch=<BASE64_VALUE>'
```

to the `curl` command above.

## Expected response

A successful response is JSON with a bearer token and expiry metadata, for
example:

```json
{
  "access_token": "<ACCESS_TOKEN>",
  "token_type": "Bearer",
  "scope": "disability_rating.read veteran_status.read",
  "expires_in": 300
}
```

Use the returned token on subsequent VA API requests:

```text
Authorization: Bearer <ACCESS_TOKEN>
```

## Postman equivalent

If you prefer Postman, VA's public `Lighthouse OAuth Token` collection already
does the JWT signing step for you. Set:

- `aud`
- `token_endpoint`
- `clientId`
- `privatePem`

Then send the `Client Credentials Example` request with the correct `scope`,
and include `launch` only when the API-specific docs require it.

## Repository note

The current Emmy veteran-verification runtime uses this JWT client-assertion
flow. Configure the following env vars for the app:

- `VA_CLIENT_ID`
- `VA_TOKEN_URL`
- `VA_AUD`
- `VA_PRIVATE_KEY_PATH`

`VA_PRIVATE_KEY_PATH` must point to a readable RSA private key PEM file on the
filesystem.

## Sources

- [VA Veteran Service History and Eligibility client-credentials docs](https://developer.va.gov/explore/api/veteran-service-history-and-eligibility/client-credentials)
- [VA Postman collection: Lighthouse OAuth Token](https://github.com/department-of-veterans-affairs/vets-api-clients/blob/master/samples/postman/Lighthouse%20OAuth%20Token.postman_collection.json)
- [VA Postman README](https://github.com/department-of-veterans-affairs/vets-api-clients/blob/master/samples/postman/README.md)
- [VA `vets-api` veteran verification client config](https://github.com/department-of-veterans-affairs/vets-api/blob/master/lib/lighthouse/veteran_verification/configuration.rb)
