# Upstream Flutter UI Contract

This repository does not ship a standalone UI. The supported UI surface is the
upstream `myhome-jarvis` Flutter Connectors tab.

## Fixture

Upstream UI may read:

```text
fixtures/public-ui-status.sample.json
```

The fixture is public-safe metadata only. It contains no raw evidence,
credentials, private archives, local paths, account identifiers, collector
configuration, or network-write permission.

## Connector Mapping

`upstream_connector` is shaped to match the upstream Flutter
`ConnectorReadiness` model:

| Fixture field | Flutter field | Notes |
| --- | --- | --- |
| `key` | `ConnectorReadiness.key` | Always `external-evidence-lake`. |
| `label` | `ConnectorReadiness.label` | Reader-facing card title. |
| `category` | `ConnectorReadiness.category` | Public boundary grouping. |
| `status` | `ConnectorReadiness.status` | Mirrors top-level `readiness`. |
| `fixture_mode` | `ConnectorReadiness.fixtureMode` | Always `true` in this repo. |
| `data_classes` | `ConnectorReadiness.dataClasses` | Public metadata classes only. |
| `allowed_operations` | `ConnectorReadiness.allowedOperations` | Read/display/link only. |
| `forbidden_operations` | `ConnectorReadiness.forbiddenOperations` | Raw import, credentials, private archives, and collector writes stay blocked. |
| `next_step` | `ConnectorReadiness.nextStep` | Safe next reader-facing action. |

## Status Mapping

- `readiness`: `bootstrap`, `ready`, or `blocked`.
- `freshness.status`: `fixture`, `fresh`, `stale`, or `unknown`.
- `validation.status`: latest public verifier posture.
- `blocked_reason`: short reader-facing reason when collection or raw evidence
  work is out of scope.
- `last_safe_check`: command and timestamp for the last public-safe verifier
  result.
- `boundaries`: booleans proving private material is absent and external writes
  are not allowed.

## Display Rules

The upstream Connectors card may show labels, states, display cards, allowed
operations, blocked operations, and last safe check metadata.

It must not show or derive:

- raw evidence payloads;
- credentials, cookies, tokens, or keys;
- private archive paths;
- local absolute paths;
- account identifiers;
- collector configuration;
- network-write controls.

## Validation

Run:

```sh
scripts/verify-public-skeleton.sh
```

The verifier checks the fixture shape, the upstream connector mapping, and the
public-safety boundary.
