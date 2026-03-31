# Codex Automations

This directory stores repository-owned automation specs for Codex.

These files are the versioned source of truth for recurring flows we want to
reuse across machines or teammates. The live scheduled automation still lives
in the user's Codex home directory rather than in this repository, so a spec in
this directory should be treated as the checked-in definition to recreate or
update the real automation.

Current specs:

- `doc-drift-repair.md`: weekly documentation drift audit and repair flow that
  works on a dedicated `codex/` branch and prepares a draft PR when safe fixes
  are available.
