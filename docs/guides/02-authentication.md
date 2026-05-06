# Authenticating with the Emmy API

## OAuth 2.0 Client Credentials

The v0 contract uses the `OAuth2ClientCredentials` security scheme. Clients
obtain an access token from the configured token endpoint and then present that
token as a bearer token on API requests.

For the checked-in v0 spec, the security scheme is defined in
`api-spec/v0/openapi.yaml` with this token URL:

```text
https://emmy-uat.auth.us-east-1.amazoncognito.com/oauth2/token
```

Your environment-specific onboarding materials may give you a different
`auth_base`; use the values provided for your environment when making real
requests.

## Token Request Example

The Postman examples in this repository use HTTP Basic auth to send the client
credentials during the client-credentials token exchange. The equivalent `curl`
request is:

```bash
curl --location '<AUTH_BASE>' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--user '<CLIENT_ID>:<CLIENT_SECRET>' \
--data-urlencode 'grant_type=client_credentials'
```

A successful response returns a bearer token:

```json
{
    "access_token": "<ACCESS_TOKEN>",
    "expires_in": 3600,
    "token_type": "Bearer"
}
```

## Using the Access Token

Once you have an access token, send it on Emmy API requests as an
`Authorization` header:

```text
Authorization: Bearer <ACCESS_TOKEN>
```

Example:

```bash
curl --location --request POST '<API_BASE>/api/v0/education-enrollments' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <ACCESS_TOKEN>' \
--data '{
    "firstName": "Lynette",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24"
}'
```
