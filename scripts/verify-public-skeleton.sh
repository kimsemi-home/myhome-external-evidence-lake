#!/usr/bin/env bash
set -euo pipefail

required_files=(
  "README.md"
  ".codex/project-goal.md"
  ".mhj/context-pack.json"
  ".mhj/hash-cache-inputs.json"
  ".github/workflows/quality.yml"
  "docs/security.md"
  "docs/private-data-policy.md"
  "docs/bootstrap-checklist.md"
  "docs/importer-boundary.md"
  "docs/upstream-flutter-ui-contract.md"
  "schemas/ui-status.schema.json"
  "fixtures/public-ui-status.sample.json"
)

for file in "${required_files[@]}"; do
  if [[ ! -s "${file}" ]]; then
    echo "required file is missing or empty: ${file}" >&2
    exit 1
  fi
done

if [[ -e "data/private" ]]; then
  echo "private data directory must not exist in this public repo" >&2
  exit 1
fi

private_identity_pattern="kimjoo""yoon|kim-joo-""yoon"
local_user_path_pattern="/""Users""/"
token_prefix_pattern="github_""pat_|gh""p_|gh""o_|gh""u_|gh""s_"
secret_shape_pattern="BEGIN (RSA|OPENSSH|DSA|EC|PRIVATE) KEY|AKIA[0-9A-Z]{16}|xox[baprs]-|sk-[A-Za-z0-9]{20,}"
scan_pattern="${private_identity_pattern}|${local_user_path_pattern}|${token_prefix_pattern}|${secret_shape_pattern}"

if grep -RInE \
  "${scan_pattern}" \
  --exclude-dir=.git \
  --exclude-dir=.mhj/cache \
  --exclude=verify-public-skeleton.sh \
  .; then
  echo "public-safety scan found forbidden private material" >&2
  exit 1
fi

jq -e '
  .upstream_repo == "kimsemi-home/myhome-jarvis" and
  .context_pack_version == "v1" and
  .ontology_version == "concept-registry/v1" and
  .security_contract_version == "security/v1" and
  .authority_contract_version == "authority/v1" and
  .private_lake_stays_private == true and
  .raw_payload_public_allowed == false and
  .external_writes_allowed == false
' .mhj/context-pack.json >/dev/null

jq -e '
  .generated_contract_verified == true and
  ([.hash_cache_inputs[].key] | sort) ==
  ([
    "context_pack_version",
    "generated_artifacts",
    "ontology_version",
    "source_descriptors",
    "workflow_dependencies"
  ] | sort) and
  all(.hash_cache_inputs[]; .public_safe == true and (.sha256 | test("^[0-9a-f]{64}$")))
' .mhj/hash-cache-inputs.json >/dev/null

jq -e '
  .schema_version == "evidence-lake-ui-status/v1" and
  .repo == "kimsemi-home/myhome-external-evidence-lake" and
  .upstream_repo == "kimsemi-home/myhome-jarvis" and
  .context_pack_version == "v1" and
  (.readiness | IN("bootstrap", "ready", "blocked")) and
  (.freshness.status | IN("fixture", "fresh", "stale", "unknown")) and
  .validation.command == "scripts/verify-public-skeleton.sh" and
  .validation.public_safety_scan == true and
  .boundaries.raw_payloads_present == false and
  .boundaries.credentials_present == false and
  .boundaries.private_archives_present == false and
  .boundaries.external_writes_allowed == false and
  (.blocked_reason | length > 0) and
  .last_safe_check.command == "scripts/verify-public-skeleton.sh" and
  (.last_safe_check.result | IN("passing", "failing", "not_run")) and
  .upstream_connector.key == "external-evidence-lake" and
  .upstream_connector.label == "External evidence lake" and
  .upstream_connector.category == "public_evidence_boundary" and
  .upstream_connector.status == .readiness and
  .upstream_connector.fixture_mode == true and
  (.upstream_connector.data_classes | index("ui_status_metadata")) and
  (.upstream_connector.allowed_operations | index("read_public_fixture")) and
  (.upstream_connector.forbidden_operations | index("raw_payload_import")) and
  (.upstream_connector.forbidden_operations | index("credential_request")) and
  (.upstream_connector.forbidden_operations | index("private_archive")) and
  (.upstream_connector.forbidden_operations | index("collector_write")) and
  (.display_cards | length >= 1) and
  all(.display_cards[]; (.id | test("^[a-z][a-z0-9-]*$")) and (.state | IN("ready", "watch", "blocked")))
' fixtures/public-ui-status.sample.json >/dev/null

echo "public skeleton verification passed"
