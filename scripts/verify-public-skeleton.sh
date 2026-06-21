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

echo "public skeleton verification passed"
