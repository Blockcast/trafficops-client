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

// apiDeliveryMethods is the API version-relative path to the /deliveryMethods API endpoint.
const (
	apiAPIDeliveryServiceXMLIDDeliveryMethods = apiTransportSessionsID + "/deliverymethods"
	apiDeliveryMethods                        = "/deliverymethods"
	apiDeliveryMethodsID                      = apiDeliveryMethods + "/%d"
)

// CreateDeliveryMethod creates the given DeliveryMethod.
func (to *Session) CreateDeliveryMethod(dm tc.DeliveryMethod, opts RequestOptions) (tc.DeliveryMethodResponse, toclientlib.ReqInf, error) {
	var reqInf toclientlib.ReqInf
	var resp tc.DeliveryMethodResponse

	if dm.TransportID == 0 && dm.Transport != "" {
		dsOpts := NewRequestOptions()
		dsOpts.QueryParameters.Set("name", dm.Transport)
		dss, _, err := to.GetTransports(dsOpts)
		if err != nil {
			err = fmt.Errorf("attempting to resolve transport name '%s' to a ID: %w", dm.Transport, err)
			return resp, reqInf, err
		}
		if len(dss.Response) == 0 {
			return resp, reqInf, fmt.Errorf("no transport with name '%s'", dm.Transport)
		}
		dm.TransportID = dss.Response[0].ID
	}
	if dm.TransportID == 0 {
		return resp, reqInf, fmt.Errorf("missing transport id")
	}

	var err error
	reqInf, err = to.post(fmt.Sprintf(apiAPIDeliveryServiceXMLIDDeliveryMethods, dm.SessionID), opts, dm, &resp)
	return resp, reqInf, err
}

// UpdateDeliveryMethod replaces the DeliveryMethod identified by ID with the one provided.
func (to *Session) UpdateDeliveryMethod(id int, dm tc.DeliveryMethod, opts RequestOptions) (tc.DeliveryMethodResponse, toclientlib.ReqInf, error) {
	dm.ID = id
	route := fmt.Sprintf("%s/%d", apiDeliveryMethods, id)
	var resp tc.DeliveryMethodResponse
	var reqInf toclientlib.ReqInf

	if dm.TransportID == 0 && dm.Transport != "" {
		dsOpts := NewRequestOptions()
		dsOpts.QueryParameters.Set("name", dm.Transport)
		dss, _, err := to.GetTransports(dsOpts)
		if err != nil {
			err = fmt.Errorf("attempting to resolve transport name '%s' to a ID: %w", dm.Transport, err)
			return resp, reqInf, err
		}
		if len(dss.Response) == 0 {
			return resp, reqInf, fmt.Errorf("no transport with name '%s'", dm.Transport)
		}
		dm.TransportID = dss.Response[0].ID
	}
	if dm.TransportID == 0 {
		return resp, reqInf, fmt.Errorf("missing transport id")
	}
	reqInf, err := to.put(route, opts, dm, &resp)
	return resp, reqInf, err
}

// GetDeliveryMethods returns all DeliveryMethods in Traffic Ops.
func (to *Session) GetDeliveryMethods(opts RequestOptions) (tc.DeliveryMethodsResponse, toclientlib.ReqInf, error) {
	var data tc.DeliveryMethodsResponse
	reqInf, err := to.get(apiDeliveryMethods, opts, &data)
	return data, reqInf, err
}

// DeleteDeliveryMethod lets you delete a DeliveryMethod. DeliveryMethods can be deleted by ID instead
// of by name if the ID is provided in the request options and the name is an
// empty string.
func (to *Session) DeleteDeliveryMethod(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.del(fmt.Sprintf(apiDeliveryMethodsID, id), opts, &alerts)
	return alerts, reqInf, err
}
