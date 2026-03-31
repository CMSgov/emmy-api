# Java JWT Token Example

A simple Java console application that demonstrates:
1. Making an HTTP POST request to obtain a JWT token
2. Using the token in a Bearer Authorization header
3. Making an authenticated HTTP POST request with JSON body
4. Printing the response to the console

## Requirements

- Java 11 or higher
- Maven 3.6+ (https://maven.apache.org/install.html)

## Running with Maven

This command compiles and runs this example, disabling any hostname verification so that you can test with self-signed SSL certificates. **Note that you would never want to do this in a Production environment.**

```bash
mvn compile exec:java -Dexec.mainClass="JwtTokenExample" -Djdk.internal.httpclient.disableHostnameVerification=true
```

## Configuration

Update the following constants in `src/main/java/JwtTokenExample.java` to match your environment:

- `TOKEN_ENDPOINT` - URL to the authentication endpoint that returns a JWT token
- `USERNAME` - Username for Basic authentication
- `PASSWORD` - Password for Basic authentication
- `API_ENDPOINT` - URL to the API endpoint that accepts the Bearer token
- `REQUEST_BODY` - JSON body to send with the second request

## Example Output

```
Step 1: Obtaining JWT token...
Token obtained successfully!
Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ...

Step 2: Calling API with Bearer token...
API Response:
{"status":"success","data":{"message":"Request processed successfully"},"id":"550e8400-e29b-41d4-a716-446655440000"}
```

## How It Works

### Step 1: Obtain JWT Token
- Encodes username and password in Base64 format
- Sends a POST request to `TOKEN_ENDPOINT` with `Authorization: Basic <encoded-credentials>` header
- Parses the JSON response and extracts the `token` field

### Step 2: Call API with Bearer Token
- Uses the JWT token from Step 1
- Sends a POST request to `API_ENDPOINT` with:
  - `Authorization: Bearer <token>` header
  - JSON content in the request body
- Prints the response to console

### Key Features
- Uses Java's built-in `HttpClient` (Java 11+) - no additional HTTP libraries needed
- Automatic Base64 encoding for Basic authentication
- JSON parsing with Gson
- Comprehensive error handling for network and HTTP errors

## Dependencies

- **Java HttpClient**: Built-in HTTP client (Java 11+) - no additional HTTP libraries needed
- **Gson 2.10.1**: For JSON parsing

Both are specified in `pom.xml` and will be automatically downloaded by Maven.

## Notes

- The application uses Java's built-in `HttpClient` for simplicity and minimal dependencies
- Error handling is included for network failures and HTTP errors
- The token is extracted from the first response and automatically used in the second request
- All configuration is done via static constants at the top of the class for easy modification
