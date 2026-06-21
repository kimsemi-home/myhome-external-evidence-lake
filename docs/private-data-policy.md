# Private Data Policy

Private evidence remains local to the upstream assistant environment. This
public repository may only contain public-safe contracts, skeleton files,
validation scripts, and generated metadata that has already been redacted.

## Allowed

- Public-safe context pack declarations.
- Hashes of public-safe upstream generated artifacts.
- Documentation describing boundaries and review gates.
- Validation scripts that reject private or local material.

## Not Allowed

- Raw evidence payloads.
- Private archives or local lake stores.
- Approval ledger contents.
- Household finance data, browser history, cookies, or account identifiers.
- Machine-specific absolute paths.

## Future Split Rule

Moving collectors, lake schemas, fixtures, or adapters into this repository
requires a separate scoped review. The current bootstrap only approves the
minimal public skeleton.
