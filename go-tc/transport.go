package tc

import (
	"time"

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

// TransportsResponse is the type of responses from Traffic Ops to GET requests
// made to its /transports API endpoint.
type TransportsResponse struct {
	Response []TransportChannel `json:"response"`
	Alerts
}

type TransportChannel struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"  required:"true"`
	LastUpdated   time.Time `json:"lastUpdated" db:"last_updated"`
	BroadcastName string    `json:"broadcastName" db:"broadcast_name"`
	Broadcast     int       `json:"broadcast" db:"broadcast"  required:"true"`
	//BondedChannels []int               `json:"bondedChannels" db:"bondedChannels"`
	BitrateKbps   int        `json:"bitrate" db:"bitrate" required:"true" description:"bitrate kbps"`
	Contention    float64    `json:"contention" db:"contention" required:"true" minimum:"1.0" default:"1.0" description:"contention ratio"`
	Modulation    Modulation `json:"modulation" db:"modulation"`
	FEC           string     `json:"fec" db:"fec" required:"true" pattern:"^(\\d+\\/\\d+)?$" example:"3/4"`
	GuardInterval string     `json:"guard_interval" db:"guard_interval" required:"true" pattern:"^(\\d+\\/\\d+)?$" example:"1/8"`

	HostName  string `json:"hostName" db:"host_name"`
	Server    int64  `json:"server" db:"server"  required:"true"`
	Interface string `json:"interface" db:"interface" required:"true"`
}

func (b TransportChannel) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.AddType(jsonschema.Object)
	schema.WithID("Transport")
	schema.WithDescription("Transport channel physical layer pipe.")

	return nil
}

type Modulation string

const (
	QPSK Modulation = "QPSK"
	PSK8 Modulation = "8PSK"

	APSK16 Modulation = "16APSK"
	APSK32 Modulation = "32APSK"

	QAM16  Modulation = "16QAM"
	QAM64  Modulation = "64QAM"
	QAM256 Modulation = "256QAM"

	NuQ16   Modulation = "NuQ16"
	NuQ64   Modulation = "NuQ64"
	NuQ256  Modulation = "NuQ256"
	NuQ1024 Modulation = "NuQ1024"
	NuQ4096 Modulation = "NuQ4096"
)

func (d Modulation) Enum() []interface{} {
	return []interface{}{QPSK, PSK8, APSK16, APSK32, QAM16, QAM64, QAM256, NuQ16, NuQ64, NuQ256, NuQ1024, NuQ4096}
}
