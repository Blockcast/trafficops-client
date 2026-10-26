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
	"net/url"

	"github.com/Blockcast/trafficops-client/go-tc"
	"github.com/Blockcast/trafficops-client/toclientlib"
)

const (
	apiSupplySiteCapacity      = "/supply/site-capacity"
	apiSupplyNetworkCommitment = "/supply/network-commitment"
)

// GetSiteCapacity reads site capacity ceilings. Filter with
// opts.QueryParameters: "cachegroup_id", and "source" (e.g. "nop-commitment").
func (to *Session) GetSiteCapacity(opts RequestOptions) (tc.SiteCapacitiesResponse, toclientlib.ReqInf, error) {
	var data tc.SiteCapacitiesResponse
	reqInf, err := to.get(apiSupplySiteCapacity, opts, &data)
	return data, reqInf, err
}

// GetNetworkCommitment reads a network's active commitment version. Filter with
// opts.QueryParameters: "network_id", "account_id"; "all=true" returns every version.
func (to *Session) GetNetworkCommitment(opts RequestOptions) (tc.NetworkCommitmentsResponse, toclientlib.ReqInf, error) {
	var data tc.NetworkCommitmentsResponse
	reqInf, err := to.get(apiSupplyNetworkCommitment, opts, &data)
	return data, reqInf, err
}

// PutNetworkCommitment stores a commitment version. Traffic Ops answers 409 when
// req.ExpectedVersion is stale or the network is owned by another source, and
// 200 with the existing version when the hash is already stored.
func (to *Session) PutNetworkCommitment(req tc.NetworkCommitmentPutRequest, opts RequestOptions) (tc.NetworkCommitmentResponse, toclientlib.ReqInf, error) {
	var data tc.NetworkCommitmentResponse
	reqInf, err := to.put(apiSupplyNetworkCommitment, opts, req, &data)
	return data, reqInf, err
}

// DeleteNetworkCommitment revokes the active version expectedVersion of a
// network's commitment and retracts its derived site capacity.
func (to *Session) DeleteNetworkCommitment(networkID, accountID, expectedVersion string, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("network_id", networkID)
	opts.QueryParameters.Set("account_id", accountID)
	opts.QueryParameters.Set("expected_version", expectedVersion)
	var data tc.Alerts
	reqInf, err := to.del(apiSupplyNetworkCommitment, opts, &data)
	return data, reqInf, err
}
