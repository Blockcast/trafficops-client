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

// apiBroadcasts is the API version-relative path to the /broadcasts API endpoint.
const apiBroadcasts = "/broadcasts"

// CreateBroadcast creates the given Broadcast.
func (to *Session) CreateBroadcast(broadcast tc.Broadcast, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if broadcast.Coordinate == 0 && broadcast.CoordinateName != "" {
		coordinateOpts := NewRequestOptions()
		coordinateOpts.QueryParameters.Set("name", broadcast.CoordinateName)
		coordinates, reqInf, err := to.GetCoordinates(coordinateOpts)
		if err != nil {
			return coordinates.Alerts, reqInf, err
		}
		if len(coordinates.Response) == 0 || coordinates.Response[0].ID == nil {
			return coordinates.Alerts, reqInf, errors.New("no coordinate with name " + broadcast.CoordinateName)
		}
		broadcast.Coordinate = *coordinates.Response[0].ID
	}

	if broadcast.TenantID <= 0 && broadcast.TenantName != "" {
		tenantOpts := NewRequestOptions()
		tenantOpts.QueryParameters.Set("name", broadcast.TenantName)
		ten, reqInf, err := to.GetTenants(tenantOpts)
		if err != nil {
			return ten.Alerts, reqInf, fmt.Errorf("attempting to resolve Tenant '%s' to an ID: %w", broadcast.TenantName, err)
		}
		if len(ten.Response) == 0 {
			return ten.Alerts, reqInf, fmt.Errorf("no Tenant named '%s'", broadcast.TenantName)
		}
		broadcast.TenantID = *ten.Response[0].ID
	}
	var alerts tc.Alerts
	reqInf, err := to.post(apiBroadcasts, opts, broadcast, &alerts)
	return alerts, reqInf, err
}

// UpdateBroadcast replaces the Broadcast identified by ID with the one provided.
func (to *Session) UpdateBroadcast(id int, broadcast tc.Broadcast, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiBroadcasts, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, broadcast, &alerts)
	return alerts, reqInf, err
}

// GetBroadcasts returns all Broadcasts in Traffic Ops.
func (to *Session) GetBroadcasts(opts RequestOptions) (tc.BroadcastsResponse, toclientlib.ReqInf, error) {
	var data tc.BroadcastsResponse
	reqInf, err := to.get(apiBroadcasts, opts, &data)
	return data, reqInf, err
}

// DeleteBroadcast lets you delete a Broadcast. Broadcasts can be deleted by ID instead
// of by name if the ID is provided in the request options and the name is an
// empty string.
func (to *Session) DeleteBroadcast(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}

	var alerts tc.Alerts
	route := fmt.Sprintf("%s/%d", apiBroadcasts, id)

	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
