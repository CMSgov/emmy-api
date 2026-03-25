# Getting Started with the Emmy API

The Emmy API is a CMS-provided secure service that enables you to connect to and retrieve information about Medicaid members in order to assist your state in determining eligibility. This is a **REST API** that sends and receives **JSON** data.

**With the Emmy API, you are able to:**

* Obtain student enrollment information
* Verify veteran disability status

## Technical Requirements & Prerequisites

You must be able to:

|Requirement|Details|
|--|--|
|Obtain Credentials|You must have worked with the CMS Emmy team to get onboarded and receive your API client ID and secret.|
|Access Emmy API Endpoints|The onboarding process will have provided you the endpoint URL. You must have outbound network/firewall access to this host.|
|Make HTTP POST (REST) Calls|Your system (or testing tool) must be capable of making HTTP POST calls where you can supply specific headers in the request.|

When you onboard with the CMS Emmy team, you will obtain credentials. You will use the values you received in this guide:

* ```client_id``` This is your client ID provided during onboarding
* ```client_secret``` The secret password for your client ID
* ```auth_base``` The endpoint used for **authentication** (but not Emmy API actions)
* ```api_base``` The endpoint used for **Emmy API actions** (but not authentication)

Later in this guide, you will also build or obtain these values as well:

* ```base64_credentials``` A Base64-encoded version of your client ID and secret
* ```access_token``` Your short-lived token used to make Emmy API calls

Let's get started!

## Using the Emmy API

To begin using the Emmy API, we will start with manual steps which ensure your credentials are valid, your connection functions, and you can construct a proper request.

### Step 1: Get a Token with your Client Credentials

To **authenticate** with the Emmy API, you will use your client ID and client secret to obtain a token for all Emmy API operations. Tokens are short-lived and must be refreshed once expired. Once you obtain a token, you will use it to make any subsequent Emmy API requests.

Obtaining a token requires that you use a **Basic Authentication** HTTP header. This means that you will join your ```<client_id>``` and ```<client_secreet>``` with a colon ```:``` and then Base64 encode the whole string. Follow the [Authentication Base64 Encoding](02-authentication.md#credential-format) instructions to build your credential string. We'll refer to the credential string as ```<base64_credentials>``` below.

For simplicity, we'll use ```curl``` to showcase getting an access token. From a terminal or command prompt, run the ```curl``` command below ([how do I install ```curl```?](../examples/v1/curl.md#installing-curl)) by substituting the appropriate variables with the values you received during onboarding:

```bash
curl --location '<AUTH_BASE>' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Authorization: Basic <BASE64_CREDENTIALS>' \
--data-urlencode 'grant_type=client_credentials'
```

After successfully running this command, a JSON response which contains your freshly minted ```<access_token>``` will appear. Here is an example (note that your response will likely not be beautifully indented):

```json
{
    "access_token":"<ACCESS_TOKEN>",
    "expires_in":3600,
    "token_type":"Bearer"
}
```

Copy the ```access_token``` value (this will be a long string of text) from the response. Now you are ready to [prepare your request payload](01-getting-started.md#step-2-prepare-your-request-payload) and call an Emmy API endpoint!

*(Review the [Authentication Guide](02-authentication.md) for more details on authenticating. Also consider browsing the [Postman Examples](../examples/v1/postman.md) or [curl Examples](../examples/v1/curl.md) for more help using these tools.)*

### Step 2: Prepare your Request Payload

Requests to the Emmy API use JSON content payloads in the request body. First, review the [Emmy API specification](https://cmsgov.github.io/emmy-api/api-spec) for the operation you would like to perform. In this example, we will make request a member's educational enrollment information from the Emmy API's ```/edu``` endpoint.

In the [API specs for ```/edu```](https://cmsgov.github.io/emmy-api/api-spec/#/Education/submitEducationVerification), you can see the structure of the request body. Build a JSON body using your member's information, like so:

```json
{
  "clientReferenceId": "tracking-id-from-your-system",
  "applicant": {
    "firstName": "Lynette",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24",
    "ssnLast4": "4321"
  }
}
```

Note that the ```"clientReferenceId"``` is a tracking value that you generate from your system side. The Emmy API does not care about nor use this value, other than to give it back to you so you can tie your requests and responses together.

With your JSON payload prepared, you can now [make a request to the Emmy API](01-getting-started.md#step-3-invoke-an-emmy-api-request).

### Step 3: Invoke an Emmy API Request

In this step, we will again use the ```curl``` tool to make a request to the Emmy API. We will continue the example from step 2 and request member information from the [Education ```/edu```](https://cmsgov.github.io/emmy-api/api-spec/#/Education/submitEducationVerification) endpoint.

This request is slightly more complex since it now has HTTP authorization, a content body, and a different endpoint. Use this example by substituting the values for your situation:

```bash
curl --location --request GET '<API_BASE>/v1/enrollment' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <ACCESS_TOKEN>' \
--data '{
    "firstName": "Lynette",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24"
}'
```
