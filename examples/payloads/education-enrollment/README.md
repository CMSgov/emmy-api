# Education Enrollment Sample JSON Payloads

The `.json` files in this folder represent sample requests and responses that one might expect while interacting with the education enrollment status API calls. The backing datasource is the National Student Clearinghouse, which has API specifications and sandbox sample data available. Useful links are as follows:

- [National Student Clearinghouse integration testing data](https://docs.studentclearinghouse.org/vs/insights-json/integration-testing)
- [OpenAPI (Swagger) Education Verification Specifications for Emmy API](https://cmsgov.github.io/emmy-api/swagger-ui/#/Education%20Verification)

## Sample Payload Information

### Sample Education Enrollment Requests

The sample **request** payloads in this folder are described below:

| File | Type | Description |
| :--- | :--- | :---------- |
| [`request-01-namedob.json`](request-01-namedob.json) | `POST` Single Request | Basic request with student name and date of birth. |
| [`request-02-namedobssn.json`](request-02-namedobssn.json) | `POST` Single Request | Basic request with student name, date of birth, and SSN. |
| [`request-03-allfields.json`](request-03-allfields.json) | `POST` Single Request | Basic request with student full name, date of birth, and SSN. |
| [`request-04-altformats1.json`](request-04-altformats1.json) | `POST` Single Request | Basic request with alternate acceptable formats for date of birth and SSN. |
| [`request-05-altformats2.json`](request-05-altformats2.json) | `POST` Single Request | Basic request with additional alternate formats for date of birth and SSN. |
| [`requestbatch-10-small1.json`](requestbatch-10-small1.json) | `POST` Batch Submission | Batch request with multiple records, submitted for asynchronous response. Submitting an asyncronous batch will provide you a synchronous response body (see [`response-10-requestbatch_body.json`](response-10-requestbatch_body.json) sample below). *Note that you must provide a unique `batchId` for each submission.* |
| [`requestbatch-11-small2.json`](requestbatch-11-small2.json) | `POST` Batch Submission | Asynchronous batch request with multiple records in a different format, without webhook callback. See [`response-10-requestbatch_body.json`](response-10-requestbatch_body.json) below for the synchronous response body sample. |
| [`requestbatch-12-large.json`](requestbatch-12-large.json) | `POST` Batch Submission | Asynchronous batch request with thousands of records. |

### Sample Education Enrollment Responses

Sample **response** payloads available in this folder are listed below. Note that these do not necessarily correlate with specific requests above, but are just representative response payloads:

| File | Type | HTTP Status | Description |
| :--- | :--- | :---------- | :---------- |
| [`response-01-found_with_no_enrollments.json`](response-01-found_with_no_enrollments.json) | Response | 200 | Knowledge of student was found, but no current enrollments. (for example, individual might have a degree, but is not currently enrolled in any school) |
| [`response-02-404_student_not_found.json`](response-02-404_student_not_found.json) | Response | 404 | Student was not found; HTTP status 404 is returned with this sample body. |
| [`response-03-enrolled_full_time.json`](response-03-enrolled_full_time.json) | Response | 200 | Student enrolled full time in a single school. |
| [`response-04-enrolled_full_time_multiple_schools.json`](response-04-enrolled_full_time_multiple_schools.json) | Response | 200 | Student enrolled full time in two schools. |
| [`response-05-enrolled_half_time.json`](response-05-enrolled_half_time.json) | Response | 200 | Student enrolled at least half time in a school. |
| [`response-06-enrolled_less_than_half_time.json`](response-06-enrolled_less_than_half_time.json) | Response | 200 | Student enrolled less than half time in a school. |
| [`response-10-requestbatch_body.json`](response-10-requestbatch_body.json) | Response | 201 | When submitting an asynchronous batch (such as [`requestbatch-10-small1.json`](requestbatch-10-small1.json) above), a synchronous response with HTTP status 201 is provided indicating success receiving the batch, as shown in this sample. |
| [`response-11-batch_status_queued.json`](response-11-batch_status_queued.json) | `GET` Response | 200 | Sample response for batch status when the batch is still queued. |
| [`response-12-batch_details_incomplete.json`](response-12-batch_details_incomplete.json) | `GET` Response | 200 | Sample representing details of the batch records when the status is not yet completed. |
| [`response-13-batch_status_completed.json`](response-13-batch_status_completed.json) | `GET` Response | 200 | Sample response when a batch's status has been completed. |
| [`response-14-batch_details_success.json`](response-14-batch_details_success.json) | `GET` Response | 200 | Sample response details payload for a completed batch. |
| [`response-15-batch_status_inprogress.json`](response-15-batch_status_inprogress.json) | `GET` Response | 200 | Sample status response for an in-progress batch. |
| [`response-16-batch_details_large_with_errors.json`](response-16-batch_details_large_with_errors.json) | `GET` Response | 200 | Sample details response for large completed batch which contained errors. |

If you have any requests or need assistance, reach out to the Emmy team at: [emmy@cms.hhs.gov](mailto:emmy@cms.hhs.gov)
