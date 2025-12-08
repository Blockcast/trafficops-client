package tc

import (
	"time"

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

// DeliveryMethodsResponse is the type of responses from Traffic Ops to GET requests
// made to its /deliverymethods API endpoint.
type DeliveryMethodsResponse struct {
	Response []DeliveryMethod `json:"response"`
	Alerts
}
type DeliveryMethodResponse struct {
	Response DeliveryMethod `json:"response"`
	Alerts
}

// DeliveryMethod is a named collection of Physical Locations within a Division.
type DeliveryMethod struct {
	ID                int       `json:"id" db:"id"`
	SessionID         int       `json:"session" db:"session"  required:"true" description:"user service session id"`
	TransportID       int       `json:"transport" db:"transport"  required:"true"`
	Transport         string    `json:"transportName" db:"transportname"`
	UserService       string    `json:"serviceId" db:"service_id"`
	UserServiceID     int       `json:"userservice" db:"userservice"`
	DeliveryService   string    `json:"xml_id" db:"xml_id"`
	DeliveryServiceID int       `json:"ds_id" db:"ds_id"`
	LastUpdated       time.Time `json:"lastUpdated" db:"last_updated"`
	models.DeliveryMethod
}

func (b DeliveryMethod) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.AddType(jsonschema.Object)
	schema.WithID("DeliveryMethod")
	schema.WithDescription("Allocation of unicast and broadcast transports for delivery of part of a user service session." +
		" The broadcast pipe must have sufficient capacity available for all instances of the associated session.")
	return nil
}
