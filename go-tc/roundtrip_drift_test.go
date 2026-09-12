package tc

import (
	"encoding/json"
	"testing"
)

// TestFieldSurvivesRoundTrip is the regression test for BLO-33496's failure
// class: a caller reads an object through this client, edits one field, and
// writes the whole object back. Any JSON key the vendored struct cannot hold is
// dropped on decode and omitted on encode, and Traffic Ops then binds SQL NULL
// over a column nobody meant to touch. There is no error and no log line.
//
// Each case below names a column that a Traffic Ops UPDATE statement actually
// writes, bound directly from the decoded request struct. Losing the key here
// means losing the column in production.
func TestFieldSurvivesRoundTrip(t *testing.T) {
	for _, c := range []struct {
		name string
		key  string
		in   string
		dst  func() any
	}{
		{"DeliveryServiceV41.extCDNEnabled", "extCDNEnabled", `{"extCDNEnabled":true}`, func() any { return &DeliveryServiceV41{} }},
		{"DeliveryServiceV41.extCDNAttributionSource", "extCDNAttributionSource", `{"extCDNAttributionSource":"cf_logpush"}`, func() any { return &DeliveryServiceV41{} }},
		{"DeliveryServiceV50.extCDNEnabled", "extCDNEnabled", `{"extCDNEnabled":true}`, func() any { return &DeliveryServiceV50{} }},
		{"DeliveryServiceV50.extCDNAttributionSource", "extCDNAttributionSource", `{"extCDNAttributionSource":"cf_logpush"}`, func() any { return &DeliveryServiceV50{} }},
		{"DeliveryServiceV50.acmeCertKeyType", "acmeCertKeyType", `{"acmeCertKeyType":"ecdsa-p384"}`, func() any { return &DeliveryServiceV50{} }},
		{"DeliveryServiceV50.acmeCertDurationDays", "acmeCertDurationDays", `{"acmeCertDurationDays":7}`, func() any { return &DeliveryServiceV50{} }},
		{"DeliveryServiceV50.acmeProfile", "acmeProfile", `{"acmeProfile":"shortlived"}`, func() any { return &DeliveryServiceV50{} }},
		{"CDNFederation.provider", "provider", `{"provider":"cloudflare"}`, func() any { return &CDNFederation{} }},
		{"CDNFederationV5.provider", "provider", `{"provider":"cloudflare"}`, func() any { return &CDNFederationV5{} }},
		{"UserServiceSession.serviceClass", "serviceClass", `{"serviceClass":"premium"}`, func() any { return &UserServiceSession{} }},
		{"UserServiceSession.transitClass", "transitClass", `{"transitClass":"bulk"}`, func() any { return &UserServiceSession{} }},
		{"UserServiceSession.footprintId", "footprintId", `{"footprintId":"fp-1"}`, func() any { return &UserServiceSession{} }},
		{"UserServiceSession.bandwidthKbps", "bandwidthKbps", `{"bandwidthKbps":4200}`, func() any { return &UserServiceSession{} }},
	} {
		t.Run(c.name, func(t *testing.T) {
			obj := c.dst()
			if err := json.Unmarshal([]byte(c.in), obj); err != nil {
				t.Fatalf("decode: %v", err)
			}
			out, err := json.Marshal(obj)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			var got map[string]json.RawMessage
			if err := json.Unmarshal(out, &got); err != nil {
				t.Fatalf("re-decode: %v", err)
			}
			if _, ok := got[c.key]; !ok {
				t.Errorf("key %q did not survive read-modify-write; Traffic Ops will bind NULL over its column", c.key)
			}
		})
	}
}
