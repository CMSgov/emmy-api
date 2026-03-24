# Documentation Accuracy Review (as of <Month DD, YYYY>)

## Section A: Executive Summary

**Report date:** `<YYYY-MM-DD>`  
**Scope:** `root docs + docs/**/*.md`  
**Overall accuracy score:** **<NN>/100**  
**Completion gate status:** **<PASS|FAIL>**

**Risk statement:**
- <Top risk 1 in plain language>
- <Top risk 2 in plain language>
- <Top risk 3 in plain language>

## Section B: Severity-Ranked Findings

| Severity | Doc location | Observed mismatch | Source-of-truth evidence on current branch | Recommended fix |
|---|---|---|---|---|
| `P1` | `<path:line>` | <Mismatch summary> | `<repo/path:line>` | <Actionable fix> |
| `P2` | `<path:line>` | <Mismatch summary> | `<repo/path:line>` | <Actionable fix> |
| `P3` | `<path:line>` | <Mismatch summary> | `<repo/path:line>` | <Actionable fix> |

## Section C: Update Backlog Checklist by Doc File

### `<docs/path-1>.md`
- [ ] <Implementation-ready update mapped to finding>
- [ ] <Implementation-ready update mapped to finding>

### `<docs/path-2>.md`
- [ ] <Implementation-ready update mapped to finding>

## Section D: Hygiene Appendix

### Commands run
- `<command 1>`
- `<command 2>`
- `<command 3>`

### Sensitive-file scan findings
- <Result summary for sensitive-file scan>

### Markdown link/reference/path validation findings
- <Result summary for markdown link/reference/path validation>

### Hygiene verdict
- **<Blocker|Non-blocker>:** <Clear decision and short justification>

## Section E: Deferred Watchlist (Non-blocking)

- If `<trigger condition>`, update:
  - `<affected/doc/path-1.md>`
  - `<affected/doc/path-2.md>`
- If `<trigger condition>`, update:
  - `<affected/doc/path-3.md>`
