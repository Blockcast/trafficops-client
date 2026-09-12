package tc

import (
	"github.com/Blockcast/multicast-api/3gpp/models"
	"github.com/swaggest/jsonschema-go"
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

// UserServiceSessionsResponse is the type of responses from Traffic Ops to GET requests
// made to its /userServices API endpoint.
type UserServiceSessionsResponse struct {
	Response []UserServiceSession `json:"response"`
	Alerts
}
type UserServiceSessionResponse struct {
	Response UserServiceSession `json:"response"`
	Alerts
}

//type  RRule.Set

type UserServiceSession struct {
	ID        int    `json:"id" db:"id"`
	ServiceId string `json:"serviceId" db:"serviceId" required:"true"`
	// ServiceClass, TransitClass, FootprintId and BandwidthKbps are written by
	// the Traffic Ops userservice_session UPDATE (service_class, transit_class,
	// footprint_id, bandwidth_kbps), bound straight from the decoded request
	// body. They must exist here or a full-object read-modify-write through
	// UpdateTransportSession silently NULLs all four. ServiceClass is nil for
	// non-reservation sessions.
	ServiceClass  *string `json:"serviceClass,omitempty" db:"service_class"`
	TransitClass  *string `json:"transitClass,omitempty" db:"transit_class"`
	FootprintId   *string `json:"footprintId,omitempty" db:"footprint_id"`
	BandwidthKbps *int    `json:"bandwidthKbps,omitempty" db:"bandwidth_kbps"`
	models.Session
	//Delivery []DeliveryMethod `json:"streams" db:"streams"`
	LastUpdated *Time `json:"lastUpdated" db:"lastUpdated"`
	//========
	//// DeliveryMethodsResponse is the type of responses from Traffic Ops to GET requests
	//// made to its /deliverymethods API endpoint.
	//type DeliveryMethodsResponse struct {
	//	Response []DeliveryMethod `json:"response"`
	//	Alerts
	//}
	//
	//// DeliveryMethod is a named collection of Physical Locations within a Division.
	//type DeliveryMethod struct {
	//	models.DeliveryMethod
	//	ID int `json:"id" db:"id"`
	//	//Name          string    `json:"name" db:"name"`
	//	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
	//	//Broadcast     int       `json:"broadcast" db:"broadcast"`
	//>>>>>>>> WIP:lib/go-tc/deliverymethod.go
}

func (b UserServiceSession) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.AddType(jsonschema.Object)
	schema.WithID("UserServiceSession")
	schema.WithDescription("Recurring user service session instances")
	return nil
}
