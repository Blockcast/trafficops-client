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
	"fmt"
	"github.com/Blockcast/trafficops-client/go-tc"
	"github.com/Blockcast/trafficops-client/toclientlib"
)

// apiTransportSessions is the API version-relative path to the /transportsessions API endpoint.
const (
	apiTransportSessions   = "/userservice_sessions"
	apiTransportSessionsID = apiTransportSessions + "/%d"
)

// CreateTransportSession creates the given UserServiceSession.
func (to *Session) CreateTransportSession(ts tc.UserServiceSession, opts RequestOptions) (tc.UserServiceSessionResponse, toclientlib.ReqInf, error) {
	var reqInf toclientlib.ReqInf
	var resp tc.UserServiceSessionResponse
	var err error
	//dsOpts := NewRequestOptions()
	//dsOpts.QueryParameters.Set("serviceId", ts.ServiceId)
	//dss, _, err := to.GetUserServices(dsOpts)
	//if err != nil {
	//	err = fmt.Errorf("attempting to resolve US name '%s' to ID: %w", ts.ServiceId, err)
	//	return resp, reqInf, err
	//}
	//if len(dss.Response) == 0 {
	//	return resp, reqInf, fmt.Errorf("no US with name '%s'", ts.ServiceId)
	//}
	//usId = int(dss.Response[0].ID)
	//
	//if usId == 0 {
	//	return resp, reqInf, fmt.Errorf("missing user service xmlId")
	//}

	reqInf, err = to.post(apiTransportSessions, opts, ts, &resp)
	return resp, reqInf, err
}

// UpdateTransportSession replaces the UserServiceSession identified by ID with the one provided.
func (to *Session) UpdateTransportSession(ts tc.UserServiceSession, opts RequestOptions) (tc.UserServiceSessionResponse, toclientlib.ReqInf, error) {
	var resp tc.UserServiceSessionResponse
	var reqInf toclientlib.ReqInf

	usId := ts.ID
	if usId == 0 {
		dsOpts := NewRequestOptions()
		dsOpts.QueryParameters.Set("serviceId", ts.ServiceId)
		dss, _, err := to.GetUserServices(dsOpts)
		if err != nil {
			err = fmt.Errorf("attempting to resolve US name '%s' to ID: %w", ts.ServiceId, err)
			return resp, reqInf, err
		}
		if len(dss.Response) == 0 {
			return resp, reqInf, fmt.Errorf("no US with name '%s'", ts.ServiceId)
		}
		usId = int(dss.Response[0].ID)
	}
	if usId == 0 {
		return resp, reqInf, fmt.Errorf("missing user service xmlId")
	}
	route := fmt.Sprintf("%s/%d", apiTransportSessions, usId)
	reqInf, err := to.put(route, opts, ts, &resp)
	return resp, reqInf, err
}

// GetTransportSessions returns all TransportSessions in Traffic Ops.
func (to *Session) GetTransportSessions(opts RequestOptions) (tc.UserServiceSessionsResponse, toclientlib.ReqInf, error) {
	var data tc.UserServiceSessionsResponse
	reqInf, err := to.get(apiTransportSessions, opts, &data)
	return data, reqInf, err
}

// DeleteTransportSession lets you delete a UserServiceSession. TransportSessions can be deleted by ID instead
// of by name if the ID is provided in the request options and the name is an
// empty string.
func (to *Session) DeleteTransportSession(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.del(fmt.Sprintf(apiTransportSessionsID, id), opts, &alerts)
	return alerts, reqInf, err
}
