# Project Goal

Build a public-safe external evidence lake companion for `myhome-jarvis`.

The repository exists to hold reusable external evidence lake scaffolding only
after an upstream authority gate allows it. It must remain inspectable by
Codex and GitHub Actions without exposing private household, finance, browser,
approval, credential, cookie, or local machine data.

## Operating Rules

- Treat `kimsemi-home/myhome-jarvis` as the upstream SSOT.
- Import context, ontology, authority, security, and version handoff from the
  upstream context pack declaration.
- Keep raw payloads and private archives outside this public repository.
- Prefer pull requests and merge evidence for all changes after bootstrap.
- Validate public safety even when cache hits allow expensive steps to skip.
