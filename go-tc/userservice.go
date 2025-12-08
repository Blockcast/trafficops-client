package tc

import (
	"github.com/Blockcast/multicast-api/3gpp/models"
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

// UserServicesResponse is the type of responses from Traffic Ops to GET requests
// made to its /userServices API endpoint.
type UserServicesResponse struct {
	Response []UserService `json:"response"`
	Alerts
}

type UserServiceResponse struct {
	Response UserService `json:"response"`
	Alerts
}
type UserService struct {
	ID int `json:"id" db:"id"`
	models.UserServiceDescription
	TenantID int     `json:"tenant_id" db:"tenant_id"`
	Tenant   *string `json:"tenant,omitempty"`
	//CDNID             int     `json:"cdnId" db:"cdn_id"`
	DeliveryServiceId int    `json:"dsId,omitempty" db:"ds_id"`
	DeliveryService   string `json:"xml_id,omitempty"`
	// CDNName is the name of the CDN to which the Delivery Service belongs.
	CDNName     *string   `json:"cdnName"`
	LastUpdated time.Time `json:"lastUpdated" db:"last_updated"`
}

type ServiceBundle struct {
	UserService
	Sessions []UserServiceSession
}
