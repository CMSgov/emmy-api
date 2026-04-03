import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Base64;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;
import java.security.cert.X509Certificate;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

/**
 * Simple Java console application demonstrating:
 * 1. HTTP POST to get JWT token using Basic authentication
 * 2. HTTP POST using JWT Bearer token with JSON body
 * 3. Printing response to console
 */
public class EmmyApiAuthAndRequestExample {

    // ============================================================
    // CONFIGURATION - Change these variables to match your needs
    // ============================================================

    // Update these variables with your actual values.
    private static final String AUTH_BASE = "https://emmy-prod.auth.us-east-1.amazoncognito.com/oauth2/token";
    private static final String CLIENT_ID = "6ja5f2dijpkc1g764tu65ontbu";
    private static final String CLIENT_SECRET = "oi06m4572eup96rtp2l8h2s48j3f4ngmvcqdl59qme8avnq3nk5";
    private static final String API_ENDPOINT = "https://api.emmy.cms.gov/api/v0/education-enrollments";

    // JSON request body for the sample API call for student enrollment.
    private static final String REQUEST_BODY_ENROLLMENT = "{\n" + //
            "    \"firstName\": \"Lynette\",\n" + //
            "    \"lastName\": \"Oyola\",\n" + //
            "    \"dateOfBirth\": \"1988-10-24\"\n" + //
            "}";

    /**
     * This is the main entry point for this example application.
     * @param args
     */
    public static void main(String[] args) {
        try {

            // Step 1: Get JWT bearer token using Basic authentication
            System.out.println("Step 1: Obtaining JWT token...");
            String token = getJwtToken();
            System.out.println("Token obtained successfully!");
            System.out.println("Token: " + token);
            System.out.println();

            // Step 2: Call an Emmy API endpoint with Bearer token
            System.out.println("Step 2: Calling Emmy API with Bearer token...");
            String response = callApiWithToken(token);
            System.out.println("API Response:");
            System.out.println(response);

        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }

    /**
     * Step 1: Make HTTP POST request to get JWT token using Basic authentication
     */
    private static String getJwtToken() throws Exception {

      // Create Basic Auth header
        String credentials = CLIENT_ID + ":" + CLIENT_SECRET;
        String encodedCredentials = Base64.getEncoder().encodeToString(credentials.getBytes());
        String basicAuthHeader = "Basic " + encodedCredentials;

        // Create HTTP request
        HttpRequest request = HttpRequest.newBuilder()
                .uri(new URI(AUTH_BASE))
                .header("Authorization", basicAuthHeader)
                .header("Content-Type", "application/x-www-form-urlencoded")
                .POST(HttpRequest.BodyPublishers.ofString("grant_type=client_credentials"))
                .build();

        // Send request with a client that accepts self-signed certificates
        // NOTE: In production, you should use a properly configured HttpClient with
        // valid SSL certificates!
        HttpClient client = createTrustAllCertificatesClient();
        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());

        // Check response status
        if (response.statusCode() != 200) {
            throw new Exception("Failed to get token. Status: " + response.statusCode() +
                    ", Response: " + response.body());
        }

        // Parse JWT token from response
        JsonObject jsonResponse = JsonParser.parseString(response.body()).getAsJsonObject();
        return jsonResponse.get("access_token").getAsString();
    }

    /**
     * Step 2: Make HTTP POST request to API endpoint with Bearer token and JSON body
     */
    private static String callApiWithToken(String token) throws Exception {

      // Create Bearer Auth header
        String bearerAuthHeader = "Bearer " + token;

        // POST the HTTP request with JSON body
        HttpRequest request = HttpRequest.newBuilder()
                .uri(new URI(API_ENDPOINT))
                .header("Authorization", bearerAuthHeader)
                .header("Content-Type", "application/json")
                .method("POST", HttpRequest.BodyPublishers.ofString(REQUEST_BODY_ENROLLMENT))
                .build();

        // Send request with a client that accepts self-signed certificates
        // NOTE: In production, you should use a properly configured HttpClient with
        // valid SSL certificates!
        HttpClient client = createTrustAllCertificatesClient();
        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());

        // Check response status
        if (response.statusCode() < 200 || response.statusCode() >= 300) {
            throw new Exception("API call failed. Status: " + response.statusCode());
        }

        return response.body();
    }

    /**
     * Create an HttpClient that accepts self-signed certificates and invalid SAN
     * WARNING: This is ONLY for development/testing. Never use in production!
     */
    private static HttpClient createTrustAllCertificatesClient() throws Exception {

      // Create a trust manager that accepts all certificates
        TrustManager[] trustAllCerts = new TrustManager[]{
            new X509TrustManager() {
                public X509Certificate[] getAcceptedIssuers() {
                    return null;
                }
                public void checkClientTrusted(X509Certificate[] certs, String authType) {}
                public void checkServerTrusted(X509Certificate[] certs, String authType) {}
            }
        };

        // Create SSL context with the trust manager
        SSLContext sslContext = SSLContext.getInstance("TLS");
        sslContext.init(null, trustAllCerts, new java.security.SecureRandom());

        // Disable hostname verification and endpoint identification
        // Note: This is a workaround for Java's HttpClient which doesn't expose
        // hostname verification settings directly in the API.
        // Note: You also need to use the Java environment setting: -Djdk.internal.httpclient.disableHostnameVerification=true
        try {
            javax.net.ssl.SSLParameters sslParams = sslContext.getDefaultSSLParameters();
            sslParams.setEndpointIdentificationAlgorithm(null);
        } catch (Exception ignore) {
            // Some SSLContext implementations don't support modification
        }

        // Create and return HttpClient with the custom SSL context
        return HttpClient.newBuilder()
                .sslContext(sslContext)
                .build();
    }
}
