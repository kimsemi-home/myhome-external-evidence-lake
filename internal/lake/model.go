package lake

type TrendContract struct {
	SchemaVersion              string   `json:"schema_version"`
	ContractID                 string   `json:"contract_id"`
	CollectorEnabled           bool     `json:"collector_enabled"`
	NetworkWritesAllowed       bool     `json:"network_writes_allowed"`
	RawPayloadsAllowed         bool     `json:"raw_payloads_allowed"`
	AccountIdentifiersAllowed  bool     `json:"account_identifiers_allowed"`
	MaxAPIDataRevalidationDays int      `json:"max_api_data_revalidation_days"`
	DecisionUse                string   `json:"decision_use"`
	SourceClasses              []string `json:"source_classes"`
	RequiredProvenance         []string `json:"required_provenance"`
	ForbiddenFields            []string `json:"forbidden_fields"`
}

type TrendObservation struct {
	SchemaVersion string      `json:"schema_version"`
	ObservationID string      `json:"observation_id"`
	Source        TrendSource `json:"source"`
	Signal        TrendSignal `json:"signal"`
	Provenance    Provenance  `json:"provenance"`
	Boundaries    Boundaries  `json:"boundaries"`
}

type TrendSource struct {
	SourceID        string `json:"source_id"`
	SourceClass     string `json:"source_class"`
	Provider        string `json:"provider"`
	Authorization   string `json:"authorization"`
	RetrievedAt     string `json:"retrieved_at"`
	ExpiresAt       string `json:"expires_at"`
	SourceURISHA256 string `json:"source_uri_sha256"`
	PayloadSHA256   string `json:"payload_sha256"`
	RightsStatus    string `json:"rights_status"`
}

type TrendSignal struct {
	Region         string `json:"region"`
	Category       string `json:"category"`
	Kind           string `json:"kind"`
	Direction      string `json:"direction"`
	ConfidenceBand string `json:"confidence_band"`
	WindowStart    string `json:"window_start"`
	WindowEnd      string `json:"window_end"`
}

type Provenance struct {
	CollectorVersion       string `json:"collector_version"`
	TransformVersion       string `json:"transform_version"`
	EventTimePreserved     bool   `json:"event_time_preserved"`
	RetrievalTimePreserved bool   `json:"retrieval_time_preserved"`
}

type Boundaries struct {
	ContainsRawPayload         bool `json:"contains_raw_payload"`
	ContainsCredentials        bool `json:"contains_credentials"`
	ContainsAccountIdentifiers bool `json:"contains_account_identifiers"`
	ContainsRevenueDetails     bool `json:"contains_revenue_details"`
	ContainsLocalPaths         bool `json:"contains_local_paths"`
	AllowsExternalWrite        bool `json:"allows_external_write"`
}
