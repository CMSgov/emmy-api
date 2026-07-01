# EnrollmentResponse V0 to V1 Transition Guide

This guide provides a mapping and migration path for consumers transitioning from Enrollment V0 to V1 APIs.

- [Request Field Mapping](#request-field-mapping)
- [Response Field Mapping](#field-mapping)

## Why are we making these changes?

To support transitioning to the hub API over time, we are introducing `EnrollmentResponseV1` to reflect the likely BSD changes that will be made in the near future. One of the perspectives of hub APIs is no inferred or guessed data.

## What are the major changes when transitioning to the hub?

- The oauth provider will be changing (moving from emmy), but the oauth mechanism will be remaining the same.
- Bulk API support will not be supported in the initial hub specs. It is being considered at a later date, but not upon initial release.
- Transaction data is exposed to support auditing.
- UAT environments for Emmy will stay up to support initial testing, but production environments will be shut down and not available for January.

## Request Field Mapping

The V1 request model (`EnrollmentRequestV1`) adopts a more descriptive naming convention and adds support for multiple names and point-in-time lookups.

| V0 Field | V1 Field | Notes                                                                                                                                                                                                           |
| :--- | :--- |:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `firstName` | `personGivenName` | Renamed.                                                                                                                                                                                                        |
| `middleName` | `personMiddleName` | Renamed.                                                                                                                                                                                                        |
| `lastName` | `personSurName` | Renamed.                                                                                                                                                                                                        |
| `dateOfBirth` | `personBirthDate` | Renamed.                                                                                                                                                                                                        |
| `ssn` | `personSocialSecurityNumber` | Renamed.                                                                                                                                                                                                        |
| `address` | *None* | V1 relies on SSN and Name matching; physical address is not currently used. Please reach out to [emmy@cms.hhs.gov](mailto:emmy@cms.hhs.gov) if this is a concern.                                               |
| *New* | `asOfDate` | Changes what the currentEnrollmentStatus is to reflect the asOfDate. Defaults to current time.                                                                                                                  |
| *New* | `previousNames` | **Purpose**: An array of objects containing `personGivenName`, `personMiddleName`, and `personSurName`. Used to search for records if the student's name has changed (e.g., marriage), increasing the hit rate. |
| *New* | `termsAcceptedIndicator` | Boolean. Explicitly confirms the student has consented to the data lookup. You must require as 'Yes'.                                                                                                           |

## Key Architectural Changes

1. **Nesting**: V0 provided a flat list of enrollment events. V1 groups enrollment events (`enrollmentData`) under their respective educational institutions (`enrollmentDetails`).
2. **Status Handling**: V0 performed "normalization" of enrollment statuses into a set of internal ranks. V1 returns the raw status codes from the data source (NSC) to provide higher fidelity.
3. **Expanded Context**: V1 includes the student information that was provided in the request and detailed transaction metadata (IDs, fees, timestamps).
4. **Error Handling**: V1 returns a `200 OK` with a success/failure indicator in `responseMetadata` and `transactionDetails` for "No Hit" scenarios, whereas V0 often resulted in a `404 Not Found` at the service level.

---

## Field Mapping

### Top-Level Fields

| V0 Field | V1 Equivalent | Notes |
| :--- | :--- | :--- |
| `enrollmentStatus` | *None* | V1 removes the aggregated top-level status. Consumers should derive this from `enrollmentDetails` if needed. |
| `dataSource` | *Implicit* | V1 currently assumes NSC but does not explicitly return this field. |
| `metadata.durationMs` | *None* | Detailed timing is currently not exposed in V1's public response. |
| *New* | `studentInfoProvided` | Echoes back the student details used for the lookup. |
| *New* | `transactionDetails` | Contains `transactionId`, `orderId`, and NSC hit indicators. |
| *New* | `responseMetadata` | Contains `responseCode` and `responseText`. |

### Enrollment Details (School Level)

V1 introduces a hierarchy. Fields that were repeated for every entry in V0 are now grouped.

| V0 Field | V1 Field | Notes                                                                                                                                                    |
| :--- | :--- |:---------------------------------------------------------------------------------------------------------------------------------------------------------|
| `schoolName` | `officialSchoolName` | Renamed for clarity.                                                                                                                                     |
| `schoolCode` | `schoolCode` | Same (6-digit OPEID).                                                                                                                                    |
| *New* | `branchCode` | two-digit indicator of the branch of the school it is (see [Appendix: Branch Code Values](https://docs.studentclearinghouse.org/vs/json-field-descriptions/appendix-branch-code-values)) |
| *New* | `currentEnrollmentStatus` | The most recent enrollment status reported by the school, as applied to the asOfDate.                                                                    |

### Enrollment Data (Term/Event Level)

In V1, these are items within the `enrollmentData` array under a school.

| V0 Field | V1 Field | Notes |
| :--- | :--- | :--- |
| `termBeginDate` | `termBeginDate` | Same (YYYY-MM-DD). |
| `termEndDate` | `termEndDate` | Same (YYYY-MM-DD). |
| `enrollmentStatus` | `enrollmentStatusCode` | **CRITICAL**: V0 returned normalized values (e.g., `FULL_TIME`). V1 returns raw source codes (e.g., `F`, `H`, `L`). |
| *New* | `schoolCertifiedOnDate` | Date the school reported the data. |
| *New* | `anticipatedGraduationDate` | Expected graduation date. |

---

## Reverse Mapping (V1 to V0)

If you need to convert a V1 response back to a V0 format for legacy compatibility, follow this logic:

### 1. Flatten the structure

Iterate through `enrollmentDetails` (schools), and for each school, iterate through its `enrollmentData` (terms).

### 2. Map fields

- `schoolName` = `officialSchoolName`
- `schoolCode` = `schoolCode`
- `termBeginDate` = `termBeginDate`
- `termEndDate` = `termEndDate`
- `enrollmentStatus` = `Normalize(enrollmentStatusCode)`

### 3. Normalize Statuses

Use the following mapping for `enrollmentStatusCode` (standard NSC codes) to convert to V0 `enrollmentStatus`:

| Raw Code | V0 Normalized Status | Description |
| :--- | :--- | :--- |
| `F` | `FULL_TIME` | Full-time |
| `Q` | `THREE_QUARTERS_TIME` | Three-quarter time |
| `H` | `HALF_TIME` | Half-time |
| `L` | `LESS_THAN_HALF_TIME` | Less than half-time |
| `Y` | `ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING` | Enrollment status unknown (default if not released by school) |

### 4. Aggregated Status (Top Level)

To recreate the V0 `enrollmentStatus`, find the "highest rank" among all normalized statuses found in the details.

If you need to find the status for a specific date (e.g., "today"), you can use logic similar to this:

#### Ruby (Model-based)

```ruby
def find_highest_status_on_date(v1_response, target_date = Date.today)
  relevant_statuses = []

  v1_response[:enrollmentDetails].each do |school|
    school[:enrollmentData].each do |term|
      begin_date = Date.parse(term[:termBeginDate])
      end_date = Date.parse(term[:termEndDate])

      if target_date.between?(begin_date, end_date)
        normalized = normalize_status(term[:enrollmentStatusCode])
        relevant_statuses << normalized if normalized
      end
    end
  end

  # Return the status with the highest rank (using the ranks defined in EnrollmentStatus)
  relevant_statuses.max_by { |status| status_rank(status) }
end
```

#### JavaScript (JSON-based)

```javascript
const RANKS = {
  'FULL_TIME': 5,
  'THREE_QUARTERS_TIME': 4,
  'HALF_TIME': 3,
  'LESS_THAN_HALF_TIME': 2,
  'ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING': 1
};

function normalizeStatus(code) {
  const mapping = {
    'F': 'FULL_TIME',
    'Q': 'THREE_QUARTERS_TIME',
    'H': 'HALF_TIME',
    'L': 'LESS_THAN_HALF_TIME',
    'Y': 'ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING'
  };
  return mapping[code] || null;
}

function findHighestStatusOnDate(v1Response, targetDateStr = new Date().toISOString().split('T')[0]) {
  const targetDate = new Date(targetDateStr);
  let bestStatus = null;

  v1Response.enrollmentDetails.forEach(school => {
    school.enrollmentData.forEach(term => {
      const beginDate = new Date(term.termBeginDate);
      const endDate = new Date(term.termEndDate);

      if (targetDate >= beginDate && targetDate <= endDate) {
        const normalized = normalizeStatus(term.enrollmentStatusCode);
        if (normalized) {
          if (!bestStatus || RANKS[normalized] > RANKS[bestStatus]) {
            bestStatus = normalized;
          }
        }
      }
    });
  });

  return bestStatus;
}
```

--- Understanding "current enrollment status"

Current enrollment status is calculated based off the asOfDate if provided, otherwise current date.

This returns CC if they are currently enrolled in that school as of the current date, or CN if not currently enrolled.

Can be used as a convenience method when searching through enrollments to only grab ones relevant for the asOfDate time period.

---

## Important Implementation Notes

- **No Hit Behavior**: In V1, there are two main ways a response might indicate no records were found, even with a `200 OK` status.
  1. **Direct No Hit**: `transactionDetails.nscHitIndicator` is `false`.
  2. **Implicit No Hit**: `transactionDetails.nscHitIndicator` is `true`, but `transactionDetails.transactionStatusCode` is `"UCF"` or the `enrollmentDetails` array is empty.

  Possible reasons for an Implicit No Hit include:
  - Enrollment records exist but degree records were not located.
  - The student may have records at a school different from the one requested.
  - The student may have blocks or other impediments preventing access to their records.

  NSC will not provide further details in these cases.

  - **Ruby Helper**: You can use the `no_hit?` helper method in the Ruby model which encapsulates both checks:

    ```ruby
    response = Education::EnrollmentResponseV1.new(params)
    if response.no_hit?
      # handle no hit
    end
    ```

  - **JSON/JavaScript**:

    ```javascript
    function isNoHit(response) {
      if (response.transactionDetails.nscHitIndicator === false) return true;
      if (response.transactionDetails.transactionStatusCode === 'UCF') return true;
      if (!response.enrollmentDetails || response.enrollmentDetails.length === 0) return true;
      return false;
    }
    ```

- **Date Formats**: Both versions use `YYYY-MM-DD`. Ensure your parsers handle potential nulls if a school hasn't provided specific dates.
