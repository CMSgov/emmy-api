# Technical Debt

## Current Limitations
- `.env.example` key casing mismatch with runtime parser expectations.
- `main.go` binds fixed `:8000` instead of config-driven port usage.
- Circuit breaker middleware currently only calls `Allow` (no success/failure feedback loop).
- Redis is hard prerequisite for several tests without automatic testcontainer setup.

## Proposed Improvements
- Align sample env file keys with actual config loader keys.
- Respect `cfg.Port` in server listen address.
- Integrate breaker transition hooks around downstream call outcomes.
- Add hermetic integration test setup for Redis.

## Implementation Steps
1. Correct docs/sample env and add config consistency test.
2. Update server bind logic and backfill startup test.
3. Refactor middleware wrapper to track handler result and call breaker hooks.
4. Add optional test profile using docker/testcontainers.

## Risks
- Breaker behavior change may alter traffic patterns under failure.
- Port binding change may affect deployment assumptions/scripts.

## Estimated Complexity
- **M**: Focused changes with meaningful runtime impact.

## Assumptions
- **High confidence:** These items are already causing friction in local development and test reliability.
