package tc

import (
	"encoding/json"
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

// BroadcastCoveragesResponse is a list of BroadcastCoverages as a response.
// swagger:response BroadcastCoveragesResponse
// in: body
type BroadcastCoveragesResponse struct {
	// in: body
	Response []BroadcastCoverage `json:"response"`
	Alerts
}

// BroadcastCoverageResponse is a single BroadcastCoverage response for Update and Create to
// depict what changed.
// swagger:response BroadcastCoverageResponse
// in: body
type BroadcastCoverageResponse struct {
	// in: body
	Response BroadcastCoverage `json:"response"`
	Alerts
}

// BroadcastCoverage is a representation of a BroadcastCoverage as it relates to the Traffic
// Ops data model.
type BroadcastCoverage struct {
	ID              int     `json:"id" db:"id"`
	Fid             int     `json:"fid" db:"fid"`
	StrongestSignal float64 `json:"strongest_signal" db:"strongest_signal" required:"true"`
	BroadcastId     int     `json:"broadcast" db:"broadcast" required:"true"`
	Broadcast       string  `json:"broadcastName" db:"broadcastName"`
	Geometry        []byte  `json:"wkb_geometry" db:"wkb_geometry" required:"true" type:"string" description:"PostGIS ST_Polygon base64 of ST_AsText"`

	// LastUpdated
	//
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

func (b BroadcastCoverage) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.AddType(jsonschema.Object)
	schema.WithID("BroadcastCoverage")
	schema.WithDescription("Broadcast signal strength topography entries based on PostGIS Polygons")
	return nil
}

// BroadcastCoverageNullable is identical to BroadcastCoverage except that its fields are
// reference values, which allows them to be nil.
type BroadcastCoverageNullable struct {

	// The BroadcastCoverage to retrieve
	//
	// ID of the BroadcastCoverage
	//
	// required: true
	ID              *int     `json:"id" db:"id"`
	Fid             *int     `json:"fid" db:"fid"`
	StrongestSignal *float64 `json:"strongest_signal" db:"strongest_signal"`
	BroadcastId     int      `json:"broadcast" db:"broadcast"`
	Broadcast       *string  `json:"broadcastName" db:"broadcastName"`
	Geometry        []byte   `json:"wkb_geometry" db:"wkb_geometry"`

	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type BroadcastCoverageFeaturesResponse struct {
	// in: body
	Response []BroadcastCoverageGeoJsonFeature `json:"response"`
	Alerts
}

type BroadcastCoverageGeoJsonFeature struct {
	Type       string          `json:"type"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties json.RawMessage `json:"properties"`
}
