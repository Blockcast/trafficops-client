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
	"errors"
	"fmt"
	"net/url"

	"github.com/Blockcast/trafficops-client/go-tc"
	"github.com/Blockcast/trafficops-client/toclientlib"
)

// apiCoverages is the API version-relative path to the /coverages API endpoint.
const apiBroadcastCoverages = "/broadcast_coverages"
const apiBroadcastCoveragesGeoJson = apiBroadcastCoverages + "/geojson"

// CreateCoverage creates the given Coverage.
func (to *Session) CreateBroadcastCoverage(coverage tc.BroadcastCoverage, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if coverage.BroadcastId == 0 && coverage.Broadcast != "" {
		broadcastOpts := NewRequestOptions()
		broadcastOpts.QueryParameters.Set("name", coverage.Broadcast)
		broadcasts, reqInf, err := to.GetBroadcasts(broadcastOpts)
		if err != nil {
			return broadcasts.Alerts, reqInf, err
		}
		if len(broadcasts.Response) == 0 {
			return broadcasts.Alerts, reqInf, errors.New("no broadcast with name " + coverage.Broadcast)
		}
		coverage.BroadcastId = broadcasts.Response[0].ID
	}

	var alerts tc.Alerts
	reqInf, err := to.post(apiBroadcastCoverages, opts, coverage, &alerts)
	return alerts, reqInf, err
}

// UpdateCoverage replaces the Coverage identified by ID with the one provided.
func (to *Session) UpdateBroadcastCoverage(id int, coverage tc.BroadcastCoverage, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiBroadcastCoverages, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, coverage, &alerts)
	return alerts, reqInf, err
}

// GetCoverages returns all Coverages in Traffic Ops.
func (to *Session) GetBroadcastCoverage(opts RequestOptions) (tc.BroadcastCoveragesResponse, toclientlib.ReqInf, error) {
	var data tc.BroadcastCoveragesResponse
	reqInf, err := to.get(apiBroadcastCoverages, opts, &data)
	return data, reqInf, err
}

// GetBroadcasts returns all Broadcasts in Traffic Ops.
func (to *Session) GetBroadcastCoverageFeatures(opts RequestOptions) (tc.BroadcastCoverageFeaturesResponse, toclientlib.ReqInf, error) {
	var data tc.BroadcastCoverageFeaturesResponse
	reqInf, err := to.get(apiBroadcastCoveragesGeoJson, opts, &data)
	return data, reqInf, err
}

// DeleteCoverage lets you delete a Coverage. Coverages can be deleted by ID instead
// of by name if the ID is provided in the request options and the name is an
// empty string.
func (to *Session) DeleteBroadcastCoverage(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}

	var alerts tc.Alerts
	route := fmt.Sprintf("%s/%d", apiBroadcastCoverages, id)

	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
