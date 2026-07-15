package lake

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var shaPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

func Verify(root string) error {
	for _, check := range []func(string) error{
		verifyRequired, verifyContextPack, verifyHashInputs, verifyUIFixture,
		verifyTrendFiles, CheckDocuments, VerifyTrace, SafetyCheck, formatCheck,
	} {
		if err := check(root); err != nil {
			return err
		}
	}
	return nil
}

func verifyRequired(root string) error {
	required := []string{
		"README.md", ".codex/project-goal.md", ".mhj/context-pack.json", ".mhj/hash-cache-inputs.json",
		".github/workflows/quality.yml", "docs/security.md", "docs/private-data-policy.md", "docs/bootstrap-checklist.md",
		"docs/importer-boundary.md", "docs/upstream-flutter-ui-contract.md", "schemas/ui-status.schema.json",
		"fixtures/public-ui-status.sample.json", "contracts/shorts-trend-evidence.json",
		"fixtures/shorts-trend-observation.sample.json", "docs-src/shorts-trend-evidence.json",
		"docs/shorts-trend-evidence.md", "traceability.json",
	}
	for _, rel := range required {
		info, err := os.Stat(filepath.Join(root, rel))
		if err != nil || info.Size() == 0 {
			return fmt.Errorf("required file is missing or empty: %s", rel)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "data", "private")); err == nil {
		return errors.New("data/private must not exist in this public repository")
	}
	return nil
}

func verifyContextPack(root string) error {
	var value struct {
		UpstreamRepo             string `json:"upstream_repo"`
		ContextPackVersion       string `json:"context_pack_version"`
		OntologyVersion          string `json:"ontology_version"`
		SecurityContractVersion  string `json:"security_contract_version"`
		AuthorityContractVersion string `json:"authority_contract_version"`
		PrivateLakeStaysPrivate  bool   `json:"private_lake_stays_private"`
		RawPayloadPublicAllowed  bool   `json:"raw_payload_public_allowed"`
		ExternalWritesAllowed    bool   `json:"external_writes_allowed"`
	}
	body, err := os.ReadFile(filepath.Join(root, ".mhj", "context-pack.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &value); err != nil {
		return err
	}
	if value.UpstreamRepo != "kimsemi-home/myhome-jarvis" || value.ContextPackVersion != "v1" ||
		value.OntologyVersion != "concept-registry/v1" || value.SecurityContractVersion != "security/v1" ||
		value.AuthorityContractVersion != "authority/v1" || !value.PrivateLakeStaysPrivate || value.RawPayloadPublicAllowed || value.ExternalWritesAllowed {
		return errors.New("context pack boundary mismatch")
	}
	return nil
}

func verifyHashInputs(root string) error {
	var value struct {
		GeneratedContractVerified bool `json:"generated_contract_verified"`
		HashCacheInputs           []struct {
			Key        string `json:"key"`
			SHA256     string `json:"sha256"`
			PublicSafe bool   `json:"public_safe"`
		} `json:"hash_cache_inputs"`
	}
	body, err := os.ReadFile(filepath.Join(root, ".mhj", "hash-cache-inputs.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &value); err != nil {
		return err
	}
	want := map[string]bool{"generated_artifacts": true, "source_descriptors": true, "workflow_dependencies": true, "context_pack_version": true, "ontology_version": true}
	if !value.GeneratedContractVerified || len(value.HashCacheInputs) != len(want) {
		return errors.New("hash cache contract mismatch")
	}
	for _, item := range value.HashCacheInputs {
		if !want[item.Key] || !item.PublicSafe || !shaPattern.MatchString(item.SHA256) {
			return fmt.Errorf("invalid hash cache input %q", item.Key)
		}
		delete(want, item.Key)
	}
	return nil
}

func verifyUIFixture(root string) error {
	var value struct {
		SchemaVersion string `json:"schema_version"`
		Repo          string `json:"repo"`
		UpstreamRepo  string `json:"upstream_repo"`
		Validation    struct {
			Command          string `json:"command"`
			PublicSafetyScan bool   `json:"public_safety_scan"`
		} `json:"validation"`
		Boundaries        Boundaries `json:"boundaries"`
		UpstreamConnector struct {
			FixtureMode         bool     `json:"fixture_mode"`
			ForbiddenOperations []string `json:"forbidden_operations"`
		} `json:"upstream_connector"`
	}
	body, err := os.ReadFile(filepath.Join(root, "fixtures", "public-ui-status.sample.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &value); err != nil {
		return err
	}
	if value.SchemaVersion != "evidence-lake-ui-status/v1" || value.Repo != "kimsemi-home/myhome-external-evidence-lake" ||
		value.UpstreamRepo != "kimsemi-home/myhome-jarvis" || value.Validation.Command != "go run ./cmd/lakectl verify" ||
		!value.Validation.PublicSafetyScan || !value.UpstreamConnector.FixtureMode || !contains(value.UpstreamConnector.ForbiddenOperations, "collector_write") {
		return errors.New("public UI fixture boundary mismatch")
	}
	return nil
}

func verifyTrendFiles(root string) error {
	var contract TrendContract
	if err := decodeStrict(filepath.Join(root, "contracts", "shorts-trend-evidence.json"), &contract); err != nil {
		return err
	}
	if err := ValidateTrendContract(contract); err != nil {
		return err
	}
	var observation TrendObservation
	if err := decodeStrict(filepath.Join(root, "fixtures", "shorts-trend-observation.sample.json"), &observation); err != nil {
		return err
	}
	return ValidateTrendObservation(observation, contract)
}

func ValidateTrendContract(contract TrendContract) error {
	if contract.SchemaVersion != "shorts.trend-contract/v1" || contract.ContractID == "" || contract.DecisionUse != "hypothesis_input_only" {
		return errors.New("trend contract identity mismatch")
	}
	if contract.CollectorEnabled || contract.NetworkWritesAllowed || contract.RawPayloadsAllowed || contract.AccountIdentifiersAllowed {
		return errors.New("public trend contract enables a forbidden capability")
	}
	if contract.MaxAPIDataRevalidationDays < 1 || contract.MaxAPIDataRevalidationDays > 30 || len(contract.SourceClasses) == 0 {
		return errors.New("trend freshness or source classes invalid")
	}
	for _, field := range []string{"raw_payload", "oauth_token", "cookie", "account_id", "channel_id", "revenue", "local_path"} {
		if !contains(contract.ForbiddenFields, field) {
			return fmt.Errorf("forbidden field %q is not declared", field)
		}
	}
	return nil
}

func ValidateTrendObservation(value TrendObservation, contract TrendContract) error {
	if value.SchemaVersion != "shorts.trend-observation/v1" || value.ObservationID == "" || !contains(contract.SourceClasses, value.Source.SourceClass) {
		return errors.New("trend observation identity or source class mismatch")
	}
	if !shaPattern.MatchString(value.Source.SourceURISHA256) || !shaPattern.MatchString(value.Source.PayloadSHA256) || value.Source.RightsStatus == "" {
		return errors.New("trend source hashes or rights status invalid")
	}
	retrieved, err := time.Parse(time.RFC3339, value.Source.RetrievedAt)
	if err != nil {
		return errors.New("invalid retrieved_at")
	}
	expires, err := time.Parse(time.RFC3339, value.Source.ExpiresAt)
	if err != nil || expires.Before(retrieved) || expires.After(retrieved.AddDate(0, 0, contract.MaxAPIDataRevalidationDays)) {
		return errors.New("trend observation exceeds revalidation window")
	}
	start, err := time.Parse(time.RFC3339, value.Signal.WindowStart)
	if err != nil {
		return errors.New("invalid signal window_start")
	}
	end, err := time.Parse(time.RFC3339, value.Signal.WindowEnd)
	if err != nil || end.Before(start) || value.Signal.Region == "" || value.Signal.Category == "" || value.Signal.Kind == "" {
		return errors.New("invalid trend signal window or dimensions")
	}
	if value.Boundaries.ContainsRawPayload || value.Boundaries.ContainsCredentials || value.Boundaries.ContainsAccountIdentifiers ||
		value.Boundaries.ContainsRevenueDetails || value.Boundaries.ContainsLocalPaths || value.Boundaries.AllowsExternalWrite {
		return errors.New("trend observation contains forbidden material")
	}
	if !value.Provenance.EventTimePreserved || !value.Provenance.RetrievalTimePreserved {
		return errors.New("event and retrieval time provenance are required")
	}
	return nil
}

func decodeStrict(path string, target any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("JSON file must contain exactly one value")
	}
	return nil
}

func formatCheck(root string) error {
	for _, base := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(root, base), func(path string, entry os.DirEntry, err error) error {
			if err != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
				return err
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			formatted, err := format.Source(body)
			if err != nil {
				return err
			}
			if !bytes.Equal(body, formatted) {
				rel, _ := filepath.Rel(root, path)
				return fmt.Errorf("unformatted Go file: %s", filepath.ToSlash(rel))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
