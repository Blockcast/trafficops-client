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
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Blockcast/trafficops-client/go-tc"
)

// commitmentJSON is a version in the shape Traffic Ops (BLO-49239) writes.
const commitmentJSON = `{"network_id":"net-1","account_id":"acct-1","source":"nop-commitment",
"version":"7d4c0c38-5b7c-4f3e-9d1a-2f0e8c1b6a11","commitment_hash":"sha256:abc","status":"active",
"valid_from":"2026-10-26T00:00:00Z",
"lines":[{"up_mbps":500,"down_mbps":100,"multicast":true,"target_asns":[13335],"target_ip_ranges":["203.0.113.0/24"],"peering":"ix"},
{"up_mbps":200,"down_mbps":0,"multicast":false}]}`

func TestSupplyClientShapes(t *testing.T) {
	var got []string // method + path + raw query + body, per request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.Path+"?"+r.URL.RawQuery+" "+string(body))
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /api/5.0/supply/site-capacity":
			io.WriteString(w, `{"response":[{"id":"r1","cachegroup_id":7,"dest_type":"asn","dest_value":"13335",
"peering":"ix","max_mbps":500,"source":"nop-commitment","commitment_hash":"sha256:abc","updated_at":"x"}]}`)
		case "GET /api/5.0/supply/network-commitment":
			io.WriteString(w, `{"response":[`+commitmentJSON+`]}`)
		case "PUT /api/5.0/supply/network-commitment":
			io.WriteString(w, `{"response":`+commitmentJSON+`}`)
		case "DELETE /api/5.0/supply/network-commitment":
			io.WriteString(w, `{"alerts":[{"text":"revoked","level":"success"}]}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	to := NewNoAuthSession(srv.URL, true, "BLO-49241-test", false, 10*time.Second)

	opts := NewRequestOptions()
	opts.QueryParameters.Set("source", "nop-commitment")
	sc, _, err := to.GetSiteCapacity(opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(sc.Response) != 1 || sc.Response[0].CommitmentHash == nil || *sc.Response[0].CommitmentHash != "sha256:abc" ||
		sc.Response[0].Peering == nil || *sc.Response[0].Peering != "ix" {
		t.Fatalf("site capacity decoded wrong: %+v", sc.Response)
	}

	nc, _, err := to.GetNetworkCommitment(RequestOptions{QueryParameters: url.Values{"network_id": {"net-1"}}})
	if err != nil {
		t.Fatal(err)
	}
	c := nc.Response[0]
	if c.Version != "7d4c0c38-5b7c-4f3e-9d1a-2f0e8c1b6a11" || c.ValidTo != nil || len(c.Lines) != 2 ||
		c.Lines[0].TargetASNs[0] != 13335 || *c.Lines[0].Peering != "ix" || c.Lines[1].Peering != nil {
		t.Fatalf("commitment decoded wrong: %+v", c)
	}

	prev := "v0"
	put, _, err := to.PutNetworkCommitment(tc.NetworkCommitmentPutRequest{NetworkCommitment: c, ExpectedVersion: &prev}, NewRequestOptions())
	if err != nil {
		t.Fatal(err)
	}
	if put.Response.CommitmentHash != "sha256:abc" {
		t.Fatalf("put response decoded wrong: %+v", put.Response)
	}

	if _, _, err := to.DeleteNetworkCommitment("net-1", "acct-1", "v1", NewRequestOptions()); err != nil {
		t.Fatal(err)
	}

	if len(got) != 4 {
		t.Fatalf("requests: %q", got)
	}
	if !strings.Contains(got[0], "source=nop-commitment") {
		t.Errorf("source filter not sent: %s", got[0])
	}
	// The PUT body is the flat commitment plus expected_version, peering on the line.
	var putBody map[string]interface{}
	if err := json.Unmarshal([]byte(got[2][strings.Index(got[2], "{"):]), &putBody); err != nil {
		t.Fatal(err)
	}
	if putBody["expected_version"] != "v0" || putBody["network_id"] != "net-1" ||
		putBody["lines"].([]interface{})[0].(map[string]interface{})["peering"] != "ix" {
		t.Errorf("PUT body wrong: %v", putBody)
	}
	for _, q := range []string{"network_id=net-1", "account_id=acct-1", "expected_version=v1"} {
		if !strings.Contains(got[3], q) {
			t.Errorf("DELETE missing %s: %s", q, got[3])
		}
	}
}
