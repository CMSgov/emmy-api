# AGENTS.md

This repository is an **Express.js (v5) + TypeScript** API/service.

This file is guidance for human/AI agents contributing to the codebase.

## Tech stack

- Runtime: Node.js
- Server: Express.js
- Language: TypeScript
- Validation/parsing: **Zod** (type-safe parsing)
- Testing: Jest
- Property-based testing: **fast-check**
- Logging: pino / pino-http

## Project layout (high level)

- `src/` — application source
  - `src/index.ts` — service entrypoint
  - `src/env.ts` — environment/config parsing
  - `src/nsc/` — NSC-related domain/API code
- `test/` — Jest tests (including fast-check properties)
- `dist/` — build output
- `nsc.md` — protocol/domain documentation

## How to run

- Dev server:
  - `pnpm dev` (runs `tsx src/index.ts`)
- Tests:
  - `pnpm test`
  - `pnpm test:watch`
- Build:
  - `pnpm build`

## Coding conventions

### Prefer functional programming patterns

- Prefer **pure functions** and explicit inputs/outputs.
- Avoid hidden global state; thread dependencies through function parameters.
- Keep effects (I/O, network, env, time) at the edges (Express handlers, clients).
- Prefer `const` and expressions over mutation.

Practical guidance:

- Put domain logic in modules that do not depend on Express.
- Express route handlers should be thin:
  1) parse/validate input
  2) call domain/service functions
  3) map result to HTTP response

### TypeScript

- Keep types close to the Zod schema source-of-truth.
- Prefer `z.infer<typeof Schema>` instead of duplicating interfaces.
- Avoid `any`; if needed, isolate and narrow ASAP.

## Zod usage (type-safe parsing)

- All external input must be validated with Zod:
  - request bodies, query params, route params
  - environment variables (`src/env.ts`)
  - responses from external services when feasible

Recommended pattern:

```ts
import { z } from "zod";

export const UserId = z.string().min(1);
export type UserId = z.infer<typeof UserId>;

export const RequestSchema = z.object({ userId: UserId });
export type Request = z.infer<typeof RequestSchema>;

export const parseRequest = (input: unknown): Request =>
  RequestSchema.parse(input);
```

- Prefer `.safeParse` when you need to control error mapping to HTTP responses.
- Convert Zod errors into consistent API errors (don’t leak internals).

If you have Zod v4 questions, use the repo’s Zod-support tooling (see existing `agents.md`).

## Express conventions

- Keep handler signatures explicit (`(req, res, next) => ...`) and avoid throwing unhandled errors.
- Centralize error handling and map known failures to appropriate status codes.
- Prefer small, composable middleware:
  - request id / logging
  - auth
  - parsing/validation

## Testing

When performing any testing task (adding/updating tests, choosing test strategy, fixing flaky tests, or writing fast-check properties), **first review**:

- `skills/javascript-testing-expert.md`

### Unit tests

- Prefer unit tests for pure functions (fast, deterministic).
- Use Jest for assertions and test running.

### Property-based tests (fast-check)

Use fast-check to validate invariants and edge cases.

Guidelines:

- Start from a property/invariant, not a single example.
- Use domain-specific arbitraries when possible.
- Keep properties deterministic (avoid real network/time).

Example structure:

```ts
import fc from "fast-check";

test("round-trip property", () => {
  fc.assert(
    fc.property(fc.string(), (s) => {
      // invariant
      expect(decode(encode(s))).toBe(s);
    })
  );
});
```

## Changes & PR checklist

- [ ] No unvalidated external inputs (use Zod)
- [ ] Handlers remain thin; domain logic is testable without Express
- [ ] New behavior has tests (unit and/or property-based)
- [ ] `pnpm test` passes
- [ ] `pnpm build` passes

## Notes

- Domain documentation lives in `nsc.md`.
- There is an existing `agents.md` with Zod-support instructions; keep it in sync with this file as needed.
