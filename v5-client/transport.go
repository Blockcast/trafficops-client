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

// apiTransports is the API version-relative path to the /transports API endpoint.
const apiTransports = "/transports"

// CreateTransport creates the given TransportChannel.
func (to *Session) CreateTransport(transport tc.TransportChannel, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if transport.Broadcast == 0 && transport.BroadcastName != "" {
		broadcastOpts := NewRequestOptions()
		broadcastOpts.QueryParameters.Set("name", transport.BroadcastName)
		broadcasts, reqInf, err := to.GetBroadcasts(broadcastOpts)
		if err != nil {
			return broadcasts.Alerts, reqInf, err
		}
		if len(broadcasts.Response) == 0 {
			return broadcasts.Alerts, reqInf, errors.New("no broadcast with name " + transport.BroadcastName)
		}
		transport.Broadcast = broadcasts.Response[0].ID
	}
	if transport.Server == 0 && transport.HostName != "" {
		serverOpts := NewRequestOptions()
		serverOpts.QueryParameters.Set("hostName", transport.HostName)
		servers, reqInf, err := to.GetServers(serverOpts)
		if err != nil {
			return servers.Alerts, reqInf, err
		}
		if len(servers.Response) != 1 {
			return servers.Alerts, reqInf, errors.New("no unique server with name " + transport.HostName)
		}
		transport.Server = int64(servers.Response[0].ID)
	}
	var alerts tc.Alerts
	reqInf, err := to.post(apiTransports, opts, transport, &alerts)
	return alerts, reqInf, err
}

// UpdateTransport replaces the TransportChannel identified by ID with the one provided.
func (to *Session) UpdateTransport(id int, transport tc.TransportChannel, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if transport.Broadcast == 0 && transport.BroadcastName != "" {
		broadcastOpts := NewRequestOptions()
		broadcastOpts.QueryParameters.Set("name", transport.BroadcastName)
		broadcasts, reqInf, err := to.GetBroadcasts(broadcastOpts)
		if err != nil {
			return broadcasts.Alerts, reqInf, err
		}
		if len(broadcasts.Response) == 0 {
			return broadcasts.Alerts, reqInf, errors.New("no broadcast with name " + transport.BroadcastName)
		}
		transport.Broadcast = broadcasts.Response[0].ID
	}
	if transport.Server == 0 && transport.HostName != "" {
		serverOpts := NewRequestOptions()
		serverOpts.QueryParameters.Set("hostName", transport.HostName)
		servers, reqInf, err := to.GetServers(serverOpts)
		if err != nil {
			return servers.Alerts, reqInf, err
		}
		if len(servers.Response) != 1 {
			return servers.Alerts, reqInf, errors.New("no unique server with name " + transport.HostName)
		}
		transport.Server = int64(servers.Response[0].ID)
	}
	route := fmt.Sprintf("%s/%d", apiTransports, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, transport, &alerts)
	return alerts, reqInf, err
}

// GetTransports returns all Transports in Traffic Ops.
func (to *Session) GetTransports(opts RequestOptions) (tc.TransportsResponse, toclientlib.ReqInf, error) {
	var data tc.TransportsResponse
	reqInf, err := to.get(apiTransports, opts, &data)
	return data, reqInf, err
}

// DeleteTransport lets you delete a TransportChannel. Transports can be deleted by ID instead
// of by name if the ID is provided in the request options and the name is an
// empty string.
func (to *Session) DeleteTransport(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	var alerts tc.Alerts
	route := fmt.Sprintf("%s/%d", apiTransports, id)
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
