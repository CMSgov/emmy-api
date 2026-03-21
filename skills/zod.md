# Skill: Authoring Zod v4 schemas (incl. codecs)

## Purpose
Help agents design, implement, and iteratively refine **Zod v4** schemas that are correct, maintainable, and aligned with v4 idioms (metadata/registries, JSON Schema conversion, recursive objects) and **codecs** (encode/decode) where appropriate.

## When to use this skill
Use when you need to:
- Model/validate API payloads, config/env, events, DB records, or user input.
- Convert between *input* and *output* representations (strings ↔ numbers/dates, payload ↔ event envelopes, etc.).
- Generate/validate JSON Schema/OpenAPI-ish representations.
- Refactor Zod 3 schemas to Zod 4 patterns.

## Zod v4 patterns to prefer
### 1) Start from “shape-first” schemas
- Prefer `z.object({ ... })`, `z.array(...)`, `z.union([...])`, `z.enum([...])`.
- Keep schemas composable: extract shared leaf schemas (e.g., `Email`, `Id`, `ISODateString`).

### 2) Be explicit about unknown keys
Zod v4’s `z.object()` strips unknown keys; JSON Schema conversion reflects this via `additionalProperties: false`.
- Use `z.object()` when you want to **strip** extras.
- Use `z.strictObject()` when you want to **error** on extras.
- Use `z.looseObject()` when you want to **allow/passthrough** extras.

### 3) Use metadata via `.meta()` / registries
Prefer `.meta()` over `.describe()` (compat only).
- Attach `id`, `title`, `description`, `examples`, deprecations, etc.
- Use registries when you need typed metadata collections.

Example:
```ts
import * as z from "zod";

export const Email = z.email().meta({
  id: "Email",
  title: "Email address",
  description: "User email address",
  examples: ["first.last@example.com"],
});
```

### 4) Json Schema conversion (v4)
- Use `z.toJSONSchema(schema)` to produce JSON Schema.
- If your schema has different input/output types (defaults, coercions, pipes/transforms), remember `io`:
  - output (default)
  - `z.toJSONSchema(schema, { io: "input" })`

### 5) Recursive objects (no casting in v4)
Use getters returning schemas for recursion:
```ts
const Category = z.object({
  name: z.string(),
  get subcategories() {
    return z.array(Category);
  },
});
```

## Codec-first design (Zod 4.1+)
Use **codecs** when you must validate *and* convert between two representations.

Typical cases:
- External API uses snake_case, internal uses camelCase.
- Event envelope ↔ payload extraction.
- Stringly-typed inputs (env vars, query params) ↔ typed outputs.
- Backward-compatible migrations.

### Codec checklist
- Define **input schema** and **output schema** separately.
- Ensure conversions are total (or fail with meaningful parse errors).
- Prefer codecs over ad-hoc `transform` when a bidirectional mapping exists.

Example (event envelope ↔ payload):
```ts
import * as z from "zod";

const Payload = z.object({ body: z.string() });

const CreateEvent = z.object({
  type: z.literal("create"),
  payload: Payload,
});

export const CreateEventCodec = z.codec(
  CreateEvent, // decoded form
  Payload,    // encoded form
  {
    to: (evt) => evt.payload,
    from: (payload) => ({ type: "create", payload }),
  },
);

// decode: Payload -> CreateEvent
// encode: CreateEvent -> Payload
```

## Process: how an agent should work
1. **Clarify the contract**
   - What is the input source? (HTTP body, query, env, DB, file)
   - Field optionality vs nullability; defaulting; unknown keys policy.
   - Any cross-field constraints.

2. **Model leaf types first**
   - IDs, timestamps, enums, constrained strings.

3. **Compose objects**
   - Use `extend/omit/pick/partial` as needed.

4. **Add refinements last**
   - Use `refine/superRefine` for cross-field validation.

5. **If conversion is needed: pick a codec**
   - If bidirectional: `z.codec`.
   - If one-way only: `transform` + `pipe` or coercion APIs.

6. **Add metadata for downstream tooling**
   - `.meta({ id, title, description, examples })`

7. **(Optional) Produce JSON Schema**
   - `z.toJSONSchema(schema, { target, io, reused, cycles, unrepresentable })`

8. **Write minimal tests**
   - `safeParse` success/failure cases; codec encode/decode round-trips.

## Using zod-mcp (when appropriate)
When you’re uncertain about an API detail, version availability, or best-practice for a particular Zod v4 feature, use the **zod-mcp server** to confirm.

Use zod-mcp to:
- Verify signatures/behavior for `z.codec`, `z.toJSONSchema`, registries/`.meta()`.
- Check edge cases (optional vs nullable, unknown keys behavior).
- Find the canonical v4 way to express a constraint (e.g., ISO formats, file, string formats).

Suggested query patterns:
- “In Zod v4, what’s the difference between `z.object`, `z.strictObject`, `z.looseObject`?”
- “How does `z.toJSONSchema` treat optional fields and `additionalProperties`?”
- “What is the exact `z.codec` signature and encode/decode direction?”

## Output expectations (what you should deliver)
- A Zod v4 schema (or schema set) with:
  - Clear naming and composition.
  - Correct optional/nullable semantics.
  - Explicit unknown-key policy.
  - Codecs where conversion is required.
  - `.meta()` fields for key public schemas.
- If requested: a JSON Schema output snippet via `z.toJSONSchema`.
