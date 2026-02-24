# Refactoring Plan

## Current Limitations
- Route handler for education mixes demo payload creation with endpoint responsibility.
- Some startup/config concerns are spread across `main` and package constructors.
- Middleware behavior documentation and tests are sparse for auth/circuit-breaker composition.

## Proposed Improvements
- Separate transport DTO handling from domain/service invocation.
- Introduce focused route-level request/response contracts.
- Improve constructor ergonomics and dependency wiring clarity.
- Expand middleware and error-path test coverage.

## Implementation Steps
1. Introduce request binding + validation for education route.
2. Move payload construction to transport layer with explicit schema.
3. Add table-driven tests for middleware chains and failure mapping.
4. Normalize startup wiring and config-driven listen address usage.
5. Re-run docs alignment after refactor completion.

## Risks
- Refactors can unintentionally alter external API behavior.
- Increased abstraction may reduce readability if over-applied.

## Estimated Complexity
- **M**: Moderate structural improvements with manageable blast radius.

## Assumptions
- **High confidence:** Refactoring now will reduce future integration churn as API surface expands.
