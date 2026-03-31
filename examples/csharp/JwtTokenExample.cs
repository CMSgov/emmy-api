using System.Net.Http.Headers;
using System.Net.Security;
using System.Text;
using System.Text.Json;

// ============================================================
// CONFIGURATION - Change these variables to match your needs
// ============================================================

// Update these variables with your actual values.
const string AuthBase = "https://us-east-1xcpcizvmn.auth.us-east-1.amazoncognito.com/oauth2/token";
const string ClientId = "7v1o4kgtd9dmkhvcm9ghpvk2vb";
const string ClientSecret = "4jo0c6aeq96fefbpnr3hu098boppe42skk2ful1ok51ofe371nc";
const string ApiEndpoint = "https://api.dev.emmy.cms.gov/enrollment";

// JSON request body for the sample API call for student enrollment.
const string RequestBody = """
    {
        "firstName": "Lynette",
        "lastName": "Oyola",
        "dateOfBirth": "1988-10-24"
    }
    """;

// ============================================================
// Main logic
// ============================================================

try
{
    // Step 1: Get JWT token using Basic authentication
    Console.WriteLine("Step 1: Obtaining JWT token...");
    string token = await GetJwtTokenAsync();
    Console.WriteLine("Token obtained successfully!");
    Console.WriteLine($"Token: {token}");
    Console.WriteLine();

    // Step 2: Call API endpoint with Bearer token
    Console.WriteLine("Step 2: Calling API with Bearer token...");
    string response = await CallApiWithTokenAsync(token);
    Console.WriteLine("API Response:");
    Console.WriteLine(response);
}
catch (Exception ex)
{
    Console.Error.WriteLine($"Error: {ex.Message}");
    Console.Error.WriteLine(ex.StackTrace);
}

// ============================================================
// Step 1: POST to token endpoint using Basic auth
// ============================================================

async Task<string> GetJwtTokenAsync()
{
    using HttpClient client = CreateTrustAllCertificatesClient();

    // Build Basic Auth header
    string credentials = $"{ClientId}:{ClientSecret}";
    string encodedCredentials = Convert.ToBase64String(Encoding.UTF8.GetBytes(credentials));
    client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Basic", encodedCredentials);

    // Send form-urlencoded body with grant_type
    var content = new FormUrlEncodedContent(new[]
    {
        new KeyValuePair<string, string>("grant_type", "client_credentials")
    });

    HttpResponseMessage response = await client.PostAsync(AuthBase, content);

    if (!response.IsSuccessStatusCode)
    {
        string body = await response.Content.ReadAsStringAsync();
        throw new Exception($"Failed to get token. Status: {(int)response.StatusCode}, Response: {body}");
    }

    string responseBody = await response.Content.ReadAsStringAsync();
    using JsonDocument doc = JsonDocument.Parse(responseBody);
    return doc.RootElement.GetProperty("access_token").GetString()
        ?? throw new Exception("access_token not found in response.");
}

// ============================================================
// Step 2: GET with JSON body using Bearer token
// ============================================================

async Task<string> CallApiWithTokenAsync(string token)
{
    using HttpClient client = CreateTrustAllCertificatesClient();

    // Set Bearer auth header
    client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);

    // GET with a body requires a custom HttpRequestMessage
    // Note: GET with a body is non-standard; .NET requires explicit use of HttpRequestMessage
    var requestMessage = new HttpRequestMessage(HttpMethod.Get, ApiEndpoint)
    {
        Content = new StringContent(RequestBody, Encoding.UTF8, "application/json")
    };

    HttpResponseMessage response = await client.SendAsync(requestMessage);

    if (!response.IsSuccessStatusCode)
    {
        throw new Exception($"API call failed. Status: {(int)response.StatusCode}");
    }

    return await response.Content.ReadAsStringAsync();
}

// ============================================================
// Helper: HttpClient that accepts self-signed / invalid SAN certs
// WARNING: ONLY for development/testing. Never use in production!
// ============================================================

HttpClient CreateTrustAllCertificatesClient()
{
    var handler = new HttpClientHandler
    {
        // Accept all certificates, including self-signed and those with invalid SANs
        ServerCertificateCustomValidationCallback = (message, cert, chain, errors) => true
    };

    return new HttpClient(handler);
}
