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

View [education response detail](./04-education-responses-detail.md) for more examples and further nuance of the API explained.

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

View [veteran disability response detail](./05-veteran-responses-detail.md) for more examples and further nuance of the API explained.
