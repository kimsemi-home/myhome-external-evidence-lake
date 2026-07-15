package lake

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestTrendObservationRejectsExpiredAndPrivateData(t *testing.T) {
	contract := validTrendContract()
	observation := validTrendObservation()
	if err := ValidateTrendObservation(observation, contract); err != nil {
		t.Fatal(err)
	}
	observation.Source.ExpiresAt = "2026-08-15T00:00:01Z"
	if err := ValidateTrendObservation(observation, contract); err == nil {
		t.Fatal("expected observation beyond thirty days to fail")
	}
	observation = validTrendObservation()
	observation.Boundaries.ContainsAccountIdentifiers = true
	if err := ValidateTrendObservation(observation, contract); err == nil {
		t.Fatal("expected account identifier boundary to fail")
	}
}

func TestStrictObservationRejectsRawPayloadField(t *testing.T) {
	body, _ := json.Marshal(validTrendObservation())
	body = append(body[:len(body)-1], []byte(`,"raw_payload":"forbidden"}`)...)
	var value TrendObservation
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err == nil {
		t.Fatal("expected raw_payload to be rejected")
	}
}

func validTrendContract() TrendContract {
	return TrendContract{
		SchemaVersion: "shorts.trend-contract/v1", ContractID: "test", MaxAPIDataRevalidationDays: 30,
		DecisionUse: "hypothesis_input_only", SourceClasses: []string{"official_platform_api"},
		ForbiddenFields: []string{"raw_payload", "oauth_token", "cookie", "account_id", "channel_id", "revenue", "local_path"},
	}
}

func validTrendObservation() TrendObservation {
	return TrendObservation{
		SchemaVersion: "shorts.trend-observation/v1", ObservationID: "observation-1",
		Source: TrendSource{
			SourceID: "source", SourceClass: "official_platform_api", Provider: "provider", Authorization: "public",
			RetrievedAt: "2026-07-15T00:00:00Z", ExpiresAt: "2026-08-14T00:00:00Z",
			SourceURISHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			PayloadSHA256:   "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", RightsStatus: "metadata",
		},
		Signal:     TrendSignal{Region: "KR", Category: "science", Kind: "momentum", WindowStart: "2026-07-14T00:00:00Z", WindowEnd: "2026-07-15T00:00:00Z"},
		Provenance: Provenance{EventTimePreserved: true, RetrievalTimePreserved: true},
	}
}
