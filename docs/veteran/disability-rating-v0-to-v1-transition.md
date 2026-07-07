# Veteran Disability Ratings V0 to V1 Transition Guide

This guide provides a mapping and migration path for consumers transitioning from Veteran Disability Ratings V0 to V1 APIs.

- [Request Field Mapping](#request-field-mapping)
- [Response Field Mapping](#response-field-mapping)

## Why are we making these changes?

To support transitioning to the hub API, we are introducing V1 to reflect standard data models and naming conventions. One of the perspectives of hub APIs is no inferred or guessed data.

## Key Architectural Changes

1. **Nesting**: Both Request and Response models now use nested structures (`vadrRequest` and `vadrResponse`) to group related fields and support extensibility.
2. **Naming Consistency**: Field names have been tweaked to reflect hub standards on naming conventions.
3. **Detailed Ratings**: V1 introduces `individualRatings` within the `sdrInformation` block, providing granular data about specific disability decisions that was not available in V0.
4. **Transparency in VA API**: the response is split into sdrInformation and ptdrInformation blocks, based off the fact that the data is split between two different VA API calls.

---

## Request Field Mapping

The V1 request model (`DisabilityRatingRequestV1`) wraps all properties inside a `vadrRequest` object and adopts a more descriptive naming convention.

| V0 Field | V1 Field (under `vadrRequest`) | Notes |
| :--- | :--- | :--- |
| `firstName` | `personGivenName` | Renamed. |
| `middleName` | `personMiddleName` | Renamed. |
| `lastName` | `personSurName` | Renamed. |
| `dateOfBirth` | `personBirthDate` | Renamed. |
| `ssn` | `personSocialSecurityNumber` | Renamed. |
| `address.street1` | `personContactInformation.streetLineOneAddress` | Grouped and renamed. |
| `address.street2` | `personContactInformation.streetLineTwoAddress` | Grouped and renamed. |
| `address.city` | `personContactInformation.cityName` | Grouped and renamed. |
| `address.state` | `personContactInformation.stateText` | Grouped and renamed. |
| `address.postalCode` | `personContactInformation.zipCode` | Grouped and renamed. |
| `address.country` | `personContactInformation.countryText` | Grouped and renamed. |
| *New* | `personSexCode` | Optional ("M", "F", "m", "f"). |

---

## Response Field Mapping

The V1 response model (`DisabilityRatingResponseV1`) wraps all properties inside a `vadrResponse` object.

### Top-Level Fields

| V0 Field | V1 Equivalent | Notes |
| :--- | :--- | :--- |
| `earliestRatingEndDate` | `vadrResponse.earliestRatingEndDate` | Remained at the top level of the response wrapper. |
| `dataSource` | *Implicit* | V1 assumes VA but does not explicitly return this field. |
| `metadata` | *None* | V0's top-level metadata is replaced by per-block `responseMetadata`. |

### PTDR Information (Permanent and Total Disability)

Grouped under `vadrResponse.ptdrInformation`.

| V0 Field | V1 Field | Notes |
| :--- | :--- | :--- |
| `permanentDisabilityStatus` | `serviceConnectedStatusIndicator` | Renamed. |
| `totalDisabilityStatus` | `totalDisabilityStatusIndicator` | Renamed. |
| `totalDisabilityStatusEffectiveDate` | `totalDisabilityEffectiveDate` | Renamed. |
| *New* | `pensionAwardStatusIndicator` | Indicates if a pension has been awarded. |
| *New* | `responseMetadata` | Contains `responseCode` and `responseText` for this block. |

### SDR Information (Service Disabled Rating)

Grouped under `vadrResponse.sdrInformation`.

| V0 Field | V1 Field | Notes                                                                                                                                  |
| :--- | :--- |:---------------------------------------------------------------------------------------------------------------------------------------|
| `combinedDisabilityRating` | `combinedDisabilityRatingPercentage` | Renamed.                                                                                                                               |
| `combinedEffectiveDate` | `combinedEffectiveDate` | Same, but moved under the sdrInformation object instead of top level.                                                                  |
| `legalEffectiveDate` | `legalEffectiveDate` | Same, but moved under the sdrInformation object instead of top level.                                                                  |
| *New* | `individualRatings` | Array of objects: `decisionText`, `ratingEffectiveDate`, `ratingEndDate`, `ratingPercentage`, `disabilityRatingId`, `staticIndicator`. |
| *New* | `responseMetadata` | Contains `responseCode` and `responseText` for this block.                                                                             |

### Understanding service connected status

For API implementation purposes, the recommended data element for determining whether VA has classified the Veteran’s total disability as permanent is `serviceConnectedStatusIndicator`.
