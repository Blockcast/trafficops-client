package tc

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

import "testing"

// TestServerV5V4ConversionsPreserveGatewayIdentity covers the second half of the
// erasure path in BLO-33496. Traffic Ops' own Update handler round-trips the
// decoded server through Downgrade() then Upgrade() before writing it, so a
// field declared on both structs but missing from either conversion is erased
// just as silently as one missing from the struct.
//
// TestServerV5DowngradeUpgrade cannot catch this: its fixture leaves these
// fields nil, so reflect.DeepEqual passes whether or not the conversions copy
// them.
func TestServerV5V4ConversionsPreserveGatewayIdentity(t *testing.T) {
	str := func(s string) *string { return &s }
	original := ServerV5{
		HardwareID:        str("ffb1e6f4-3d05-4b2a-9f0c-2a7c1d8e5b33"),
		NetworkID:         str("blockcast"),
		CacheSoftwareType: str("ATS"),
	}

	v4 := original.Downgrade()
	if v4.HardwareID == nil || v4.NetworkID == nil || v4.CacheSoftwareType == nil {
		t.Fatalf("Downgrade dropped fields: hardwareId=%v networkId=%v cacheSoftwareType=%v",
			v4.HardwareID, v4.NetworkID, v4.CacheSoftwareType)
	}

	back := v4.Upgrade()
	for _, c := range []struct {
		name string
		got  *string
		want *string
	}{
		{"hardwareId", back.HardwareID, original.HardwareID},
		{"networkId", back.NetworkID, original.NetworkID},
		{"cacheSoftwareType", back.CacheSoftwareType, original.CacheSoftwareType},
	} {
		if c.got == nil || *c.got != *c.want {
			t.Errorf("%s not preserved through Downgrade/Upgrade: got %v, want %q", c.name, c.got, *c.want)
		}
	}
}
