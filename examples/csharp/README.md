# C# JWT Token Example

A simple .NET console application that demonstrates:
1. Making an HTTP POST request to obtain a JWT token using Basic authentication
2. Using the token in a Bearer Authorization header
3. Making an authenticated GET request with a JSON body
4. Printing the response to the console

## Requirements

- .NET 10.0 SDK or higher

## Configuration

Update the following constants in the section marked **CONFIGURATION** in the `EmmyApiAuthAndRequestExample.cs` file to match your environment:

- `AuthBase` - URL to the authentication endpoint that returns a JWT token
- `ClientId` - Client ID for Basic authentication
- `ClientSecret` - Client secret for Basic authentication
- `ApiEndpoint` - URL to the API endpoint that accepts the Bearer token
- `RequestBody` - JSON body to send with the API request

## Running

```bash
dotnet run
```

## Example Output

```
Step 1: Obtaining JWT token...
Token obtained successfully!
Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...

Step 2: Calling API with Bearer token...
API Response:
{"status":"success","data":{"message":"Request processed successfully"}}
```

## How It Works

### Step 1: Obtain JWT Token
- Encodes `ClientId:ClientSecret` in Base64
- Sends a POST to `AuthBase` with `Authorization: Basic <encoded>` header
- Sends `grant_type=client_credentials` as a form-urlencoded body (standard OAuth 2.0 client credentials flow)
- Parses the JSON response and extracts the `access_token` field

### Step 2: Call API with Bearer Token
- Uses the JWT token from Step 1
- Sends a GET request to `ApiEndpoint` with:
  - `Authorization: Bearer <token>` header
  - JSON content in the request body (uses `HttpRequestMessage` since `HttpClient.GetAsync` doesn't support a body)
- Prints the response to the console

### Key Features
- Uses .NET's built-in `HttpClient` and `System.Text.Json` — no additional libraries needed
- Accepts self-signed certificates and invalid SANs via `ServerCertificateCustomValidationCallback`
- Top-level statements keep the code minimal and readable

## Notes

> **⚠️ Warning**: The `CreateTrustAllCertificatesClient()` helper disables all certificate validation.
> This is intentional for development/testing against environments with self-signed certificates
> or mismatched SANs. **Never use this in production.**
