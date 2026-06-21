# myhome-external-evidence-lake

Public-safe external evidence lake boundary for `kimsemi-home/myhome-jarvis`.

This repository starts as a minimal skeleton. The upstream SSOT, authority
gate, ontology, security policy, and version handoff stay in
`kimsemi-home/myhome-jarvis` until a later approved split moves concrete lake
code here.

## Current Scope

- Declare the downstream context pack for external evidence work.
- Preserve a public-safe bootstrap checklist and private data boundary.
- Run public-safety checks on every pull request and main push.
- Cache bootstrap units by hash so unchanged work avoids unnecessary expensive
  steps while generated contracts are still validated.

## Explicit Non-Scope

- No raw evidence payloads.
- No credentials, cookies, tokens, local paths, or private approval ledgers.
- No private local lake archive.
- No external collection adapters until a separate scoped approval exists.

## Upstream Contract

- Upstream repo: `kimsemi-home/myhome-jarvis`
- Context pack version: `v1`
- Ontology version: `concept-registry/v1`
- Authority contract: `authority/v1`
- Security contract: `security/v1`
