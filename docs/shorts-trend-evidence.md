# Public Shorts Trend Evidence Contract

The public repository carries only provenance-rich trend metadata; raw source payloads and operating identities remain private.

## Allowed evidence

- Official platform API metadata.
- Primary scientific source descriptors.
- Official economic indicator descriptors.
- Public news and community metadata with provenance and rights state.
- Hash-addressed transform receipts and bounded confidence bands.

## Forbidden evidence

- Raw payload bodies, captions, comments, transcripts, or media.
- OAuth material, cookies, browser sessions, account IDs, or channel IDs.
- Private revenue, household, or local filesystem details.
- Collector configuration that enables network writes.

## Freshness and decision use

Non-authorized YouTube API data expires no later than thirty days after retrieval and must then be refreshed or deleted. Trend observations are hypothesis inputs only; they cannot directly approve claims, publication, or scaling.

- Keep event time separate from retrieval time.
- Bind source URI and payload hashes without publishing either body.
- Preserve support and contradiction observations.
- Let the private evidence gate combine this metadata with claim, rights, originality, disclosure, and policy evidence.

## Verification

```sh
go run ./cmd/lakectl docs check
go run ./cmd/lakectl verify
go test ./...
```
