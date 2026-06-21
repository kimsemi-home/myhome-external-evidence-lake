# Security Policy

This repository is public. It must remain safe to clone, index, cache, and run
in GitHub Actions.

## Public-Safety Rules

- Do not commit credentials, cookies, API tokens, SSH keys, private keys, or
  browser/session material.
- Do not commit local absolute paths or machine-specific user names.
- Do not commit raw evidence payloads, approval ledgers, finance records, or
  household-private material.
- Do not add collection adapters that write outside this repository without a
  new scoped authority approval.
- Keep GitHub Actions read-only unless a later reviewed change grants a narrower
  permission.

## Required Checks

Every pull request and main push must run:

- public-safety scan
- context pack contract validation
- hash-cache input validation
- bootstrap checklist validation
