# Veteran Disability Sample JSON Payloads

The `.json` files in this folder represent sample requests and responses that one might expect while interacting with the veteran disability rating API calls. The backing datasource is the Veterans Affairs developer portal, which has API specifications and sandbox sample data available. Useful links are as follows:

- [VA Developer Test Data page](https://developer.va.gov/explore/api/veteran-service-history-and-eligibility/test-users)
- [OpenAPI (Swagger) Veteran Disability Rating Specifications for Emmy API](https://cmsgov.github.io/emmy-api/swagger-ui/#/Veteran%20Verification/getVeteranDisabilityStatus)

## Sample Payload Information

### Sample Veteran Disability Requests

The sample **request** payloads in this folder are described below:

| File | Type | Description |
| :--- | :--- | :---------- |
| [`request-01-namedobssn.json`](request-01-namedobssn.json) | Single Request | Basic request with veteran's name, date of birth, and SSN. |
| [`request-02-namedobaddress.json`](request-02-namedobaddress.json) | Single Request | Basic request with veteran's name, date of birth, and address (without SSN). |

### Sample Veteran Disability Responses

Sample **response** payloads available in this folder are listed below. Note that these do not necessarily correlate with specific requests above, but are just representative response payloads:

| File | Type | HTTP Status | Description |
| :--- | :--- | :---------- | :---------- |
| [`response-01-multiple_ratings.json`](response-01-multiple_ratings.json) | Response | 200 | Veteran payload with a combined disability rating below 100, including multiple individual ratings. |
| [`response-02-total_disability_with_under_100.json`](response-02-total_disability_with_under_100.json) | Response | 200 | Veteran with a total disability status of true, but a combined disability rating less than 100%. |
| [`response-03-404_veteran_not_found.json`](response-03-404_veteran_not_found.json) | Response | 404 | Response body when a veteran is not found. |

If you have any requests or need assistance, reach out to the Emmy team at: [emmy@cms.hhs.gov](mailto:emmy@cms.hhs.gov)
