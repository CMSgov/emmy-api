# Getting Started with the Emmy API

The Emmy API is a CMS-provided secure service that enables you to connect to and retrieve information about Medicaid members in order to assist your state in determining eligibility. This is a **REST API** that sends and receives **JSON** data.

**With the Emmy API, you are able to:**

- Obtain student enrollment information
- Verify veteran disability status

## Technical Requirements & Prerequisites

You must be able to:

| Requirement                 | Details                                                                                                                                                                         |
| --------------------------- |---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Obtain Credentials          | Reach out to [emmy@cms.hhs.gov](emmy@cms.hhs.gov), requesting a sandbox credential and working with the CMS emmy team to get onboarded and receive your client ID and secret via encrypted channel. |
| Access Emmy API Endpoints   | The onboarding process will have provided you the endpoint URL. You must have outbound network/firewall access to this host.                                                    |
| Make HTTP POST (REST) Calls | Your system (or testing tool) must be capable of making HTTP POST calls where you can supply specific headers in the request.                                                   |

When you onboard with the CMS Emmy team, you will obtain credentials. You will use the values you received in this guide:

- `client_id` This is your client ID provided during onboarding
- `client_secret` The secret password for your client ID
- `auth_base` The OAuth 2.0 token endpoint for your environment
- `api_base` The base URL for Emmy API actions in your environment

Later in this guide, you will obtain:

- `access_token` Your short-lived bearer token used to make Emmy API calls

Let's get started!

## Using the Emmy API

To begin using the Emmy API, we will start with manual steps which ensure your credentials are valid, your connection functions, and you can construct a proper request.

### Step 1: Get a Token with your Client Credentials

The v0 contract uses OAuth 2.0 client credentials. Obtain a bearer token from
the token endpoint for your environment before calling Emmy API operations.

For simplicity, we'll use `curl` to showcase getting an access token. The
Postman examples in this repository use HTTP Basic auth to send the client
credentials, so this guide does the same:

```bash
curl --location '<AUTH_BASE>' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--user '<CLIENT_ID>:<CLIENT_SECRET>' \
--data-urlencode 'grant_type=client_credentials'
```

After successfully running this command, a JSON response which contains your freshly minted `<access_token>` will appear. Here is an example (note that your response will likely not be beautifully indented):

```json
{
    "access_token": "<ACCESS_TOKEN>",
    "expires_in": 3600,
    "token_type": "Bearer"
}
```

Copy the `access_token` value from the response. Now you are ready to
[prepare your request payload](01-getting-started.md#step-2-prepare-your-request-payload)
and call an Emmy API endpoint.

_(Review the [Authentication Guide](02-authentication.md) for more details on authenticating. Also consider browsing the [Postman Examples](../examples/v1/postman.md) or [curl Examples](../examples/v1/curl.md) for more help using these tools.)_

### Step 2: Prepare your Request Payload

The checked-in v0 contract defines two public operations:

- `POST /api/v0/education-enrollments`
- `POST /api/v0/veteran-disability-ratings`

Both operations reuse the same request body shape from
`schema/v0/identity.schema.json`. Build a JSON body using the person's
identifying information, like so:

```json
{
    "firstName": "Lynette",
    "middleName": "Marie",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24",
    "ssn": "123-45-6789"
}
```

Required fields are `firstName`, `lastName`, and `dateOfBirth`. `middleName`
and `ssn` are optional in the current v0 schema. Veteran disability requests
may also include an optional `address` object for demographic matching when
SSN data is not available, for example:

```json
{
    "firstName": "Lynette",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24",
    "address": {
        "street1": "17020 Tortoise St",
        "city": "Round Rock",
        "state": "TX",
        "postalCode": "78664",
        "country": "USA"
    }
}
```

With your JSON payload prepared, you can now
[make a request to the Emmy API](01-getting-started.md#step-3-invoke-an-emmy-api-request).

### Step 3: Invoke an Emmy API Request

In this step, we will use `curl` to call the education verification operation.
The v0 contract expects an HTTP `POST` request with a bearer token and JSON
request body.

Use this example by substituting the values for your situation:

```bash
curl --location --request POST '<API_BASE>/api/v0/education-enrollments' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <ACCESS_TOKEN>' \
--data '{
    "firstName": "Lynette",
    "middleName": "Marie",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24",
    "ssn": "123-45-6789"
}'
```

The current v0 success response for education returns:

```json
{
    "enrollmentStatus": "FULL_TIME"
}
```

To request veteran disability data instead, send the same identity payload to:

```text
POST <API_BASE>/api/v0/veteran-disability-ratings
```
