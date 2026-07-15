# Importer Boundary

This repository may import only public-safe context pack metadata from
`kimsemi-home/myhome-jarvis`. It must not import raw evidence, household data,
credentials, private approval ledgers, private lake archives, browser/session
material, or machine-specific local paths.

## Allowed Fields

- upstream repo name
- context pack version
- ontology, authority, and security contract versions
- public-safe generated artifact hashes
- validation status
- readiness, freshness, and blocked-state metadata for upstream UI cards

## Rejected Fields

- raw evidence payloads
- credentials, cookies, tokens, or keys
- account identifiers or household finance records
- local absolute paths or machine usernames
- private archive locations
- collector configuration that performs network writes

## UI Metadata

Upstream Flutter UI should read fixture-shaped status from
`fixtures/public-ui-status.sample.json` and treat it as public-safe status only.
The schema is `schemas/ui-status.schema.json`.

The status card contract is intentionally small:

- `readiness`: bootstrap, ready, or blocked
- `freshness`: public fixture/freshness state
- `validation`: command and public-safety scan status
- `last_safe_check`: verifier command, timestamp, and result
- `boundaries`: booleans proving absent private material
- `upstream_connector`: the exact public-safe connector card payload shape for
  `myhome-jarvis`
- `display_cards`: short reader-facing cards for upstream shadcn-style UI

## Validation

Run:

```sh
go run ./cmd/lakectl verify
```

The verifier checks required files, scans for private material, validates the
context pack, validates hash-cache inputs, and checks the public UI status
fixture shape.
