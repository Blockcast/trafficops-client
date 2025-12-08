package tc

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

import (
	"time"

	multicast "github.com/Blockcast/multicast-api"
)

// DeliveryReportsResponse is the type of responses from Traffic Ops to GET requests
// made to its /deliveryreports API endpoint.
type DeliveryReportsResponse struct {
	Response []DeliveryReport `json:"response"`
	Alerts
}

// DeliveryReportResponse is the type of responses from Traffic Ops to POST requests
// made to its /deliveryreports API endpoint.
type DeliveryReportResponse struct {
	Response *DeliveryReport `json:"response,omitempty"`
	Alerts
}

// DeleteResponse is the type of responses from Traffic Ops to DELETE requests.
type DeleteResponse struct {
	Alerts
}

// DeliveryReport represents a delivery report for multicast sessions
type DeliveryReport struct {
	// Embed the multicast reception report
	multicast.BlockcastReceptionReport `json:",inline"`

	// Traffic Ops specific fields
	ID          int        `json:"id" db:"id"`
	TenantID    int        `json:"tenant_id" db:"tenant_id"`
	Tenant      *string    `json:"tenant,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated" db:"last_updated"`
}
