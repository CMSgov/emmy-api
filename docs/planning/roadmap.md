# Roadmap

## Current Limitations
- API surface is minimal (`/`, `/status`, `/api/edu`).
- `/api/edu` uses hardcoded request input.
- Limited test isolation from external dependencies.

## Proposed Improvements
- Introduce versioned API contract for verification requests.
- Add request validation and explicit error schema.
- Expand provider adapters beyond NSC as needed.

## Implementation Steps
1. Define target API contract and versioning strategy.
2. Implement request binding/validation for education endpoint.
3. Add API-level integration tests with deterministic fixtures.
4. Add endpoint-level metrics and SLO dashboarding.
5. Document migration path for external clients.

## Risks
- Contract changes may break early consumers.
- Upstream provider variability may destabilize test fixtures.
- Auth policy tightening may impact existing integrations.

## Estimated Complexity
- **M**: Moderate cross-package changes with API behavior impact.

## Assumptions
- **Medium confidence:** Product direction prioritizes formal external API contract in near-term milestones.
