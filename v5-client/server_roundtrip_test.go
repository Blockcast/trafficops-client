package client

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Blockcast/trafficops-client/go-tc"
)

// fakeTO is the minimum Traffic Ops stand-in needed to observe a full-object
// read-modify-write.
//
// Its row is held as a raw JSON object rather than a tc.ServerV5 on purpose:
// Traffic Ops stores columns, not this client's Go struct, so a field the
// vendored struct has drifted away from still exists server-side. Keeping the
// row untyped is also what lets this test compile — and fail — against a client
// that is missing the fields entirely, which is the whole point of a drift
// regression test.
//
// The PUT handler deliberately mirrors TO's updateQuery
// (traffic_ops/traffic_ops_golang/server/servers.go): every column named in the
// UPDATE is written from the decoded request body unconditionally, so a key the
// client omits is bound as SQL NULL rather than left alone. Do not make this
// handler defensive — it only reproduces the defect if it stays as unforgiving
// as the real one.
type fakeTO struct {
	row map[string]interface{}
}

// updatedColumns are the nullable columns TO's updateQuery writes
// unconditionally. `cache_software_type` is deliberately absent: it is not in
// updateQuery, so it survives a full-object write even when the client drops it.
var updatedColumns = []string{"hardwareId", "networkId", "interfaces"}

func (f *fakeTO) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"response": []interface{}{f.row},
			})
		case http.MethodPut:
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decoding PUT body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			for _, col := range updatedColumns {
				// A key absent from the body decodes to nil, which NamedQuery
				// binds as SQL NULL. This is the erasure.
				f.row[col] = body[col]
			}
			json.NewEncoder(w).Encode(tc.Alerts{})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (f *fakeTO) column(t *testing.T, name string) string {
	t.Helper()
	v, ok := f.row[name]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("column %q is %T, want string", name, v)
	}
	return s
}

const (
	testHardwareID = "ffb1e6f4-3d05-4b2a-9f0c-2a7c1d8e5b33"
	testNetworkID  = "blockcast"
)

// TestServerUpdateRoundTripPreservesGatewayIdentity is the regression test for
// BLO-33496. A caller that reads a server, mutates one unrelated field and
// writes the whole object back must not erase columns it never touched. This is
// exactly the shape of Blockcast/magma's beacon public_address_indexer, which
// updates only Interfaces[0].IPAddresses but PUTs the full object.
//
// Against a vendored go-tc whose ServerV5 lacks HardwareID/NetworkID this fails:
// encoding/json silently discards both keys on the GET, the PUT body omits them,
// and TO binds SQL NULL over the orc8r gateway identity — no error, no log line.
func TestServerUpdateRoundTripPreservesGatewayIdentity(t *testing.T) {
	const serverID = 42
	to := &fakeTO{row: map[string]interface{}{
		"id":         serverID,
		"hostName":   "traffic-monitor",
		"domainName": "gw-74284180.example",
		"hardwareId": testHardwareID,
		"networkId":  testNetworkID,
		"interfaces": []interface{}{map[string]interface{}{
			"name": "eth0",
			"ipAddresses": []interface{}{map[string]interface{}{
				"address":        "192.0.2.10",
				"serviceAddress": true,
			}},
		}},
	}}

	srv := httptest.NewServer(to.handler(t))
	defer srv.Close()

	session := NewNoAuthSession(srv.URL, true, "BLO-33496-test", false, 10*time.Second)

	opts := NewRequestOptions()
	opts.QueryParameters = url.Values{"id": []string{"42"}}
	got, _, err := session.GetServers(opts)
	if err != nil {
		t.Fatalf("GetServers: %v", err)
	}
	if len(got.Response) != 1 {
		t.Fatalf("GetServers returned %d servers, want 1", len(got.Response))
	}

	// Mutate one unrelated field, exactly as the magma beacon indexer does.
	server := got.Response[0]
	if len(server.Interfaces) != 1 {
		t.Fatalf("GetServers returned %d interfaces, want 1", len(server.Interfaces))
	}
	server.Interfaces[0].IPAddresses = []tc.ServerIPAddress{{
		Address:        "198.51.100.7",
		ServiceAddress: true,
	}}
	if _, _, err := session.UpdateServer(server.ID, server, NewRequestOptions()); err != nil {
		t.Fatalf("UpdateServer: %v", err)
	}

	if got := to.column(t, "hardwareId"); got != testHardwareID {
		t.Errorf("hardware_id erased by an unrelated update: got %q, want %q", got, testHardwareID)
	}
	if got := to.column(t, "networkId"); got != testNetworkID {
		t.Errorf("network_id erased by an unrelated update: got %q, want %q", got, testNetworkID)
	}

	// The change the caller actually intended must still have landed, otherwise
	// the assertions above could pass on a client that writes nothing at all.
	after, _, err := session.GetServers(opts)
	if err != nil {
		t.Fatalf("GetServers after update: %v", err)
	}
	if len(after.Response) != 1 || len(after.Response[0].Interfaces) != 1 ||
		len(after.Response[0].Interfaces[0].IPAddresses) != 1 ||
		after.Response[0].Interfaces[0].IPAddresses[0].Address != "198.51.100.7" {
		t.Errorf("intended interface update did not land: %+v", after.Response)
	}
}
