package tc

import "time"

// SupplyNodeStake is the DB model for supply_node_stake.
// Staking overlay keyed on hardware_id (orc8r hardware UUID).
// A single hardware_id maps to multiple server rows via server.hardware_id
// (relay, moq-relay, amt-relay, TM are all separate servers on one gateway).
// Operator identity (name, tenant, location) lives in NOP, not here.
type SupplyNodeStake struct {
	ID             string    `json:"id" db:"id"`
	HardwareID     string    `json:"hardware_id" db:"hardware_id"`
	StakeTxHash    string    `json:"stake_tx_hash" db:"stake_tx_hash"`
	StakeAmount    float64   `json:"stake_amount" db:"stake_amount"`
	StakedAt       time.Time `json:"staked_at" db:"staked_at"`
	ChallengeKey   string    `json:"-" db:"challenge_key"`
	OnboardingStep int       `json:"onboarding_step" db:"onboarding_step"`
	RegisteredAt   time.Time `json:"registered_at" db:"registered_at"`
}

// StakeNodeRequest is POST /supply/nodes/stake body.
// Called by the cloud service after validating operator identity against NOP.
type StakeNodeRequest struct {
	HardwareID   string  `json:"hardware_id"`
	StakeTxHash  string  `json:"stake_tx_hash"`
	StakeAmount  float64 `json:"stake_amount"`
	ChallengeKey string  `json:"challenge_key"`
}

// StakeNodeResponse is returned after successful staking.
type StakeNodeResponse struct {
	ID         string `json:"id"`
	HardwareID string `json:"hardware_id"`
}

// NodeStorageRequest is PUT /supply/nodes/{id}/storage body.
type NodeStorageRequest struct {
	Disks []DiskClaim `json:"disks"`
	RAMMB int         `json:"ram_mb"`
}

// DiskClaim is a single disk allocation.
type DiskClaim struct {
	Location  string `json:"location"`
	SizeGB    int    `json:"size_gb"`
	IsVirtual bool   `json:"is_virtual"`
}

// InterfaceCommitment is the operator's bandwidth pledge for one network
// interface on one server. CDN cache traffic is overwhelmingly egress, so a
// single MaxMbpsOut number is enough; we don't track ingress separately.
// Multicast capability is per-interface (separate from supply_multicast_reach,
// which tracks the destinations this server can reach).
//
// Wire format mirrors the columns the handler writes onto TO's existing
// `interface` table (max_bandwidth in kbps; multicast bool added by migration
// 2026041700000001).
type InterfaceCommitment struct {
	InterfaceName string `json:"interface_name"`
	MaxMbpsOut    int    `json:"max_mbps_out"`
	Multicast     bool   `json:"multicast"`
}

// SiteCapacity is a site-level peering/transit capacity ceiling.
type SiteCapacity struct {
	ID           string  `json:"id" db:"id"`
	CachegroupID int64   `json:"cachegroup_id" db:"cachegroup_id"`
	DestType     string  `json:"dest_type" db:"dest_type"`
	DestValue    *string `json:"dest_value,omitempty" db:"dest_value"`
	Peering      *string `json:"peering,omitempty" db:"peering"`
	MaxMbps      int     `json:"max_mbps" db:"max_mbps"`
	Source       string  `json:"source" db:"source"`
	// CommitmentHash is the hash of the NetworkCommitment version a
	// source=nop-commitment row was derived from; nil for other sources.
	CommitmentHash *string `json:"commitment_hash,omitempty" db:"commitment_hash"`
	UpdatedAt      string  `json:"updated_at" db:"updated_at"`
}

// SiteCapacitiesResponse is the GET /supply/site-capacity response.
type SiteCapacitiesResponse struct {
	Response []SiteCapacity `json:"response"`
	Alerts
}

// NetworkCommitment is one version of a network's capacity commitment
// (supply_network_commitment). NOP computes CommitmentHash; Traffic Ops stores
// it verbatim and uses it as the idempotency key.
type NetworkCommitment struct {
	NetworkID      string                  `json:"network_id" db:"network_id"`
	AccountID      string                  `json:"account_id" db:"account_id"`
	Source         string                  `json:"source" db:"source"`
	Version        string                  `json:"version" db:"version"`
	CommitmentHash string                  `json:"commitment_hash" db:"commitment_hash"`
	Status         string                  `json:"status" db:"status"`
	ValidFrom      time.Time               `json:"valid_from" db:"valid_from"`
	ValidTo        *time.Time              `json:"valid_to,omitempty" db:"valid_to"`
	Lines          []NetworkCommitmentLine `json:"lines" db:"lines"`
}

// NetworkCommitmentLine is one line of a commitment. Each target becomes one
// derived site-capacity row carrying the line's Peering class; a line with no
// targets derives dest_type any.
type NetworkCommitmentLine struct {
	UpMbps         int      `json:"up_mbps"`
	DownMbps       int      `json:"down_mbps"`
	Multicast      bool     `json:"multicast"`
	TargetASNs     []int64  `json:"target_asns,omitempty"`
	TargetIPRanges []string `json:"target_ip_ranges,omitempty"`
	// Peering is ix, peered, transit, or nil (unclassified).
	Peering *string `json:"peering,omitempty"`
}

// NetworkCommitmentPutRequest is the PUT /supply/network-commitment body.
// ExpectedVersion is the compare-and-set guard: nil for a first commitment,
// otherwise the active version the caller read (stale -> 409).
type NetworkCommitmentPutRequest struct {
	NetworkCommitment
	ExpectedVersion *string `json:"expected_version,omitempty"`
}

// NetworkCommitmentResponse is the PUT /supply/network-commitment response.
type NetworkCommitmentResponse struct {
	Response NetworkCommitment `json:"response"`
	Alerts
}

// NetworkCommitmentsResponse is the GET /supply/network-commitment response:
// the active version, or every version with ?all=true.
type NetworkCommitmentsResponse struct {
	Response []NetworkCommitment `json:"response"`
	Alerts
}

// FederationUsage is a federation cost attribution log entry.
type FederationUsage struct {
	ID               string   `json:"id" db:"id"`
	FederationTarget string   `json:"federation_target" db:"federation_target"`
	DSID             int64    `json:"ds_id" db:"ds_id"`
	ClientASN        *int     `json:"client_asn,omitempty" db:"client_asn"`
	BytesServed      int64    `json:"bytes_served" db:"bytes_served"`
	CostPerGB        *float64 `json:"cost_per_gb,omitempty" db:"cost_per_gb"`
	PeriodStart      string   `json:"period_start" db:"period_start"`
	PeriodEnd        string   `json:"period_end" db:"period_end"`
	LoggedAt         string   `json:"logged_at" db:"logged_at"`
}

// MulticastReach is a multicast reach claim per server.
type MulticastReach struct {
	ID         string  `json:"id" db:"id"`
	ServerID   int64   `json:"server_id" db:"server_id"`
	ReachType  string  `json:"reach_type" db:"reach_type"`
	ReachValue string  `json:"reach_value" db:"reach_value"`
	Verified   bool    `json:"verified" db:"verified"`
	VerifiedAt *string `json:"verified_at,omitempty" db:"verified_at"`
	ClaimedAt  string  `json:"claimed_at" db:"claimed_at"`
}

// BandwidthCommitmentRequest is PUT /supply/nodes/{id}/bandwidth body.
// Each entry sets the operator's per-interface upload pledge + multicast flag
// directly on TO's interface table. PoB skip-or-challenge gating happens via
// supply_node_stake.onboarding_step at the gateway level, not per-interface.
type BandwidthCommitmentRequest struct {
	Interfaces []InterfaceCommitment `json:"interfaces"`
}

// NodeSummary is a row in the My Nodes table.
// NodeName is omitted — NOP owns it, cloud service enriches.
type NodeSummary struct {
	ID             string `json:"id"`
	HardwareID     string `json:"hardware_id"`
	CacheGroup     string `json:"cache_group"`
	IPAddress      string `json:"ip_address"`
	Status         string `json:"status"`
	OnboardingStep int    `json:"onboarding_step"`
}

// PlacementRequest is the body for POST /supply/nodes/{id}/placement[/preview].
// Override preserves admin intent across deterministic re-runs.
type PlacementRequest struct {
	Override bool `json:"override"`
}

// PlacementCachegroup is the cachegroup selected by placement assignment.
type PlacementCachegroup struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

// PlacementInputs is the auditable input set captured for assignment re-runs.
type PlacementInputs struct {
	IP              string            `json:"ip"`
	ASN             int               `json:"asn"`
	GeoIP           PlacementGeoIP    `json:"geoip"`
	LatLongSelf     *PlacementLatLong `json:"latlong_self,omitempty"`
	RTT             []string          `json:"rtt"`
	DatasetVersions map[string]string `json:"dataset_versions"`
	Flags           []string          `json:"flags"`
	Override        bool              `json:"override"`
	Version         int               `json:"version"`
}

// PlacementGeoIP is the server-observed IP geolocation used by assignment.
type PlacementGeoIP struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Src string  `json:"src"`
}

// PlacementLatLong is the operator self-reported location, used only as a cross-check.
type PlacementLatLong struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// PlacementResponse is returned by placement preview and commit endpoints.
type PlacementResponse struct {
	Cachegroup   PlacementCachegroup `json:"cachegroup"`
	Market       string              `json:"market"`
	NetworkClass string              `json:"network_class"`
	Status       string              `json:"status"`
	Inputs       PlacementInputs     `json:"inputs"`
}

// PlacementShadowResponse is returned after a shadow-only placement batch run.
type PlacementShadowResponse struct {
	Population   int                   `json:"population"`
	Distribution map[string]int        `json:"distribution"`
	StatusCounts map[string]int        `json:"status_counts"`
	Totals       PlacementShadowTotals `json:"totals"`
	// DeferReasons breaks the aggregate deferred_dataset_unavailable count down
	// by the specific lookup/feed that missed (e.g. "asn_unresolved",
	// "asn_not_in_peeringdb", "geo_unresolved", "rpki_stale", "v6_*"). It
	// distinguishes "datasets not delivered" from "this node's IP/ASN isn't in a
	// present dataset" without source-diving. Empty when nothing deferred.
	DeferReasons map[string]int `json:"defer_reasons,omitempty"`
	// DeferDetails lists one entry per deferred node with the selected IPs and
	// the reason, so operators can see exactly why each node deferred.
	DeferDetails []PlacementDeferDetail `json:"defer_details,omitempty"`
	// DerivedCachegroups counts the cachegroup_supply_meta rows touched by an
	// opt-in derive pass (?derive_cachegroups=true), grouped by network_class
	// (ix/dc/res). It is the production "make cachegroup_supply_meta non-zero"
	// signal: rows are derived only from the supply zones the placement path
	// produced for real gateways, never by bulk-marking existing CDN
	// cachegroups. Nil/omitted unless the derive pass ran.
	DerivedCachegroups map[string]int `json:"derived_cachegroups,omitempty"`
}

// PlacementDeferDetail explains why a single node deferred in a shadow run.
type PlacementDeferDetail struct {
	NodeID string `json:"node_id"`
	IPv4   string `json:"ipv4,omitempty"`
	IPv6   string `json:"ipv6,omitempty"`
	Reason string `json:"reason"`
}

// PlacementShadowTotals summarizes the shadow placement batch by terminal outcome.
type PlacementShadowTotals struct {
	Population                 int `json:"population"`
	Assigned                   int `json:"assigned"`
	DeferredDatasetUnavailable int `json:"deferred_dataset_unavailable"`
	ManualReview               int `json:"manual_review"`
	NoIP                       int `json:"no_ip"`
}

// CachegroupSupplyResponse is GET /cachegroups/{id}/supply.
type CachegroupSupplyResponse struct {
	CachegroupID       int64                     `json:"cachegroup_id"`
	Members            []NodeSummary             `json:"members"`
	Market             string                    `json:"market"`
	NetworkClass       string                    `json:"network_class"`
	Status             string                    `json:"status"`
	ActivationProgress *SupplyActivationProgress `json:"activation_progress,omitempty"`
	Capacity           []SiteCapacity            `json:"capacity"`
}

// SupplyActivationProgress exposes pending-zone activation gate progress for
// onboarding explainers. Values are sampled on the membership/config-sync batch
// clock from Traffic Ops server status, which is populated by Traffic Monitor.
type SupplyActivationProgress struct {
	MembersHealthy         int     `json:"members_healthy"`
	RequiredMembersHealthy int     `json:"required_members_healthy"`
	DaysElapsed            float64 `json:"days_elapsed"`
	RequiredDaysHealthy    int     `json:"required_days_healthy"`
	VerifiedMbps           int     `json:"verified_mbps"`
	RequiredVerifiedMbps   int     `json:"required_verified_mbps"`
	Eligible               bool    `json:"eligible"`
	HealthHistorySource    string  `json:"health_history_source"`
}
