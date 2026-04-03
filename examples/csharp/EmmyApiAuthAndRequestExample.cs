using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;

// ============================================================
// CONFIGURATION - Change these variables to match your needs
// ============================================================

// Update these variables with your actual values.
const string AuthBase = "https://onboarding_authentication_endpoint";
const string ClientId = "your_client_id_here";
const string ClientSecret = "your_client_secret_here";
const string ApiEndpoint = "https://onboarding_emmy_api_endpoint/api/v0/education-enrollments";

// JSON request body for the sample API call for student enrollment.
const string RequestBody_Enrollment = """
    {
        "firstName": "Lynette",
        "lastName": "Oyola",
        "dateOfBirth": "1988-10-24"
    }
    """;


// ============================================================
// This is the main entry point for this example application.
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
    Console.WriteLine("Step 2: Calling Emmy API with Bearer token...");
    string response = await CallApiWithTokenAsync(token);
    Console.WriteLine("API Response:");
    Console.WriteLine(response);
}
catch (Exception ex)
{
    Console.Error.WriteLine($"Error: {ex.Message}");
    Console.Error.WriteLine(ex.StackTrace);
}

/// <summary>
/// Step 1: Make HTTP POST request to get JWT token using Basic authentication
/// </summary>
async Task<string> GetJwtTokenAsync()
{
    // Send request with a client that accepts self-signed certificates
    // NOTE: In production, you should use a properly configured HttpClient with
    // valid SSL certificates!
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

/// <summary>
/// Step 2: Make HTTP POST request to API endpoint with Bearer token and JSON body
/// </summary>
async Task<string> CallApiWithTokenAsync(string token)
{
    // Send request with a client that accepts self-signed certificates
    // NOTE: In production, you should use a properly configured HttpClient with
    // valid SSL certificates!
    using HttpClient client = CreateTrustAllCertificatesClient();

    // Set Bearer auth header
    client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);

    // POST the HTTP request with JSON body
    var requestMessage = new HttpRequestMessage(HttpMethod.Post, ApiEndpoint)
    {
        Content = new StringContent(RequestBody_Enrollment, Encoding.UTF8, "application/json")
    };

    HttpResponseMessage response = await client.SendAsync(requestMessage);

    if (!response.IsSuccessStatusCode)
    {
        throw new Exception($"API call failed. Status: {(int)response.StatusCode}");
    }

    return await response.Content.ReadAsStringAsync();
}

/// <summary>
/// Creates an HttpClient that accepts all SSL certificates, including self-signed and those with invalid SANs.
/// WARNING: This should ONLY be used for development and testing purposes. Do NOT use this in production environments.
/// </summary>
HttpClient CreateTrustAllCertificatesClient()
{
    var handler = new HttpClientHandler
    {
        // Accept all certificates, including self-signed and those with invalid SANs
        ServerCertificateCustomValidationCallback = (message, cert, chain, errors) => true
    };

    return new HttpClient(handler);
}
