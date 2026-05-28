# Emmy API Typical Usage Examples

These examples mirror the checked-in v0 public contract.

## Education Verification

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

Example success response:

```json
{
    "enrollmentStatus": "FULL_TIME"
}
```

View [education response detail](./04-education-responses-detail.md) for more examples and further nuance of the APi explained.

## Veteran Disability Verification

```bash
curl --location --request POST '<API_BASE>/api/v0/veteran-disability-ratings' \
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

Example success response:

```json
{
    "combinedDisabilityRating": 70
}
```

View [veteran disability response detail](./05-veteran-responses-detail.md) for more examples and further nuance of the APi explained.

## Batch API Verification

When working with more than 1000 records in a minute, we recommend using the batch API.

The batch API requires a unique batch ID.

```bash
export guid="$(uuidgen)"

curl -k --location 'https://api.uat.emmy.cms.gov/api/v0/batch-education-enrollments' \
  --header 'Content-Type: application/json' \
  --header "Authorization: Bearer ${JWT}" \
  --data "{
    \"batchId\": \"${guid}\",
    \"submittedBy\": \"org-unit-42\",
    \"callbackUrl\": \"https://your-system.example.com/webhooks/enrollment\",
    \"students\": [
      {
        \"recordId\": \"rec-001\",
        \"firstName\": \"Stacey\",
        \"lastName\": \"Peacock\",
        \"dateOfBirth\": \"1996-05-24\"
      },
      {
              \"recordId\": \"rec-002\",
              \"firstName\": \"Fakey\",
              \"lastName\": \"McFakeface\",
              \"dateOfBirth\": \"1991-01-01\"
            }
    ]
  }"
```

You can then request the status of the batch using the batch ID.

```bash
curl -k --location "https://api.uat.emmy.cms.gov/api/v0/batch-education-enrollments/${guid}" \
--header 'Content-Type: application/json' \
--header "Authorization: Bearer ${JWT}"


curl -k --location "https://api.uat.emmy.cms.gov/api/v0/batch-education-enrollments/${guid}/details" \
--header 'Content-Type: application/json' \
--header "Authorization: Bearer ${JWT}"
```
