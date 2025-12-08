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
	"github.com/Blockcast/trafficops-client/go-tc"
	"github.com/Blockcast/trafficops-client/toclientlib"
	"net/http"
)

// apiUserServices is the API version-relative path to the /userservices API endpoint.
const (
	apiUserServices   = "/userservices"
	apiUserServicesID = apiUserServices + "/%d"
)

// CreateUserService creates the given UserService.
func (to *Session) CreateUserService(userservice tc.UserService, opts RequestOptions) (tc.UserServiceResponse, toclientlib.ReqInf, error) {
	var resp tc.UserServiceResponse
	if userservice.DeliveryServiceId == 0 && userservice.DeliveryService != "" {
		opts := RequestOptions{QueryParameters: map[string][]string{"xml_id": {userservice.DeliveryService}}}
		ds, reqInf, err := to.GetDeliveryServices(opts)
		if err != nil {
			return resp, reqInf, err
		} else if len(ds.Response) == 0 {
			reqInf.StatusCode = http.StatusBadRequest
			return resp, reqInf, err
		}
		userservice.DeliveryServiceId = *ds.Response[0].ID
	}
	if userservice.DeliveryServiceId == 0 {
		var reqInf toclientlib.ReqInf
		reqInf.StatusCode = http.StatusBadRequest
		return resp, reqInf, errors.New("missing delivery service id")
	}
	reqInf, err := to.post(apiUserServices, opts, userservice, &resp)
	return resp, reqInf, err
}

// UpdateUserService replaces the UserService identified by ID with the one provided.
func (to *Session) UpdateUserService(id int, userservice tc.UserService, opts RequestOptions) (tc.UserServiceResponse, toclientlib.ReqInf, error) {
	var resp tc.UserServiceResponse
	if userservice.DeliveryServiceId == 0 && userservice.DeliveryService != "" {
		opts := RequestOptions{QueryParameters: map[string][]string{"xml_id": {userservice.DeliveryService}}}
		ds, reqInf, err := to.GetDeliveryServices(opts)
		if err != nil {
			return resp, reqInf, err
		} else if len(ds.Response) == 0 {
			reqInf.StatusCode = http.StatusBadRequest
			return resp, reqInf, err
		}
		userservice.DeliveryServiceId = *ds.Response[0].ID
	}
	if userservice.DeliveryServiceId == 0 {
		var reqInf toclientlib.ReqInf
		reqInf.StatusCode = http.StatusBadRequest
		return resp, reqInf, errors.New("missing delivery service id")
	}
	route := fmt.Sprintf(apiUserServicesID, id)
	var data tc.UserServiceResponse
	reqInf, err := to.put(route, opts, userservice, &data)
	return data, reqInf, err
}

// GetUserServices returns all UserServices in Traffic Ops.
func (to *Session) GetUserServices(opts RequestOptions) (tc.UserServicesResponse, toclientlib.ReqInf, error) {
	var data tc.UserServicesResponse
	reqInf, err := to.get(apiUserServices, opts, &data)
	return data, reqInf, err
}

// DeleteUserService lets you delete a UserService. UserServices can be deleted by ID instead
// of by name if the ID is provided in the request options and the name is an
// empty string.
func (to *Session) DeleteUserService(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.del(fmt.Sprintf(apiUserServicesID, id), opts, &alerts)
	return alerts, reqInf, err
}
