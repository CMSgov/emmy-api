# Product Questions

## Q1
- **Question:** What minimum verification domains must be supported after education (income, employment, residency, etc.)?
- **Context:** Current code only implements education verification provider flow.
- **Why it matters:** Domain roadmap determines API shape, abstraction strategy, and staffing priorities.
- **Suggested investigation direction:** Create prioritized capability matrix with stakeholder input.
- **Owner:** TBD
- **Target decision date:** TBD

## Q2
- **Question:** Who are the first external consumers, and what compatibility guarantees are required?
- **Context:** API contract is still evolving and partially scaffolded.
- **Why it matters:** Versioning, deprecation policy, and release governance depend on consumer commitments.
- **Suggested investigation direction:** Define consumer onboarding model and API maturity tiers.
- **Owner:** TBD
- **Target decision date:** TBD

## Q3
- **Question:** What compliance and data retention rules apply to verification payload/response fields?
- **Context:** Education requests may include PII-like identifiers and personal data.
- **Why it matters:** Logging, storage, caching, and redaction policies must satisfy compliance obligations.
- **Suggested investigation direction:** Capture policy requirements and map to technical controls.
- **Owner:** TBD
- **Target decision date:** TBD
