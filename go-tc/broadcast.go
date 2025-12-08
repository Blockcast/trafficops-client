package tc

import (
	"github.com/blockcast/multicast/common"
	"time"
)

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// BroadcastsResponse is the type of responses from Traffic Ops to GET requests
// made to its /broadcasts API endpoint.
type BroadcastsResponse struct {
	Response []Broadcast `json:"response"`
	Alerts
}

// A Broadcast is a named collection of Physical Locations within a Division.
type Broadcast struct {
	Latitude            float64             `json:"latitude" db:"latitude"`
	Longitude           float64             `json:"longitude" db:"longitude"`
	CoordinateName      string              `json:"coordinateName" db:"coordinate_name"`
	Coordinate          int                 `json:"coordinate" db:"coordinate"`
	ID                  int                 `json:"id" db:"id"`
	LastUpdated         time.Time           `json:"lastUpdated" db:"last_updated"`
	Name                string              `json:"name" db:"name"`
	FrequencyMhz        int                 `json:"frequency" db:"frequency"`
	Polarity            *string             `json:"polarity,omitempty" db:"polarity"`
	BandwidthMhz        int                 `json:"bandwidth" db:"bandwidth"`
	Access              common.DeliveryMode `json:"access" db:"access"`
	EffectivePowerWatts int                 `json:"erp" db:"erp"`
	TenantID            int                 `json:"tenant_id" db:"tenant"`
	Public              bool                `json:"public" db:"public"`
	TenantName          string              `json:"tenantName,omitempty"`
}

// BroadcastName is a response to a request to get a broadcast by its name. It
// includes the division that the broadcast is in.
type BroadcastName struct {
	ID       int                     `json:"id"`
	Name     string                  `json:"name"`
	Division BroadcastNameCoordinate `json:"coordinate"`
}

// BroadcastNameDivision is the division that contains the broadcast that a request
// is trying to query by name.
type BroadcastNameCoordinate struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// BroadcastNameResponse models the structure of a response to a request to get a
// broadcast by its name.
type BroadcastNameResponse struct {
	Response []BroadcastName `json:"response"`
}
