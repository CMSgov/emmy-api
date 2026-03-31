# Emmy API Typical Usage Examples

These examples mirror the checked-in v0 public contract.

## Education Verification

```bash
curl --location --request POST '<API_BASE>/v0/education-enrollments' \
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

## Veteran Disability Verification

```bash
curl --location --request POST '<API_BASE>/v0/veteran-disability-ratings' \
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
