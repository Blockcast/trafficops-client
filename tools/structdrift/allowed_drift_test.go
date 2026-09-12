package main

import (
	"os"
	"strings"
	"testing"
)

// TestEveryAllowedEntryHasAReason enforces the one rule allowed-drift.txt states
// about itself: an entry with no reason is not accepted.
//
// This exists because of an Ally finding on trafficops-client#3. The drift check
// itself cannot run in this repository -- Blockcast/trafficcontrol is private and
// this repository is public, so it runs from trafficcontrol instead (BLO-33506).
// That leaves a gap the finding named: a pull request here can add a baseline
// exemption, and nothing would object until the next scheduled run elsewhere.
//
// This does not prove a reason is TRUE; nothing mechanical can. It makes adding
// an exemption require writing one, which puts the claim in the diff where a
// reviewer has to read it. That converts a silent weakening into a visible one,
// which is what the finding actually asked for.
//
// The other half of the gap is covered by go-tc/roundtrip_drift_test.go, which
// pins the fields measured to be erasing and does not consult this file at all,
// so no edit here can weaken it.
func TestEveryAllowedEntryHasAReason(t *testing.T) {
	b, err := os.ReadFile("allowed-drift.txt")
	if err != nil {
		t.Fatalf("reading allowed-drift.txt: %v", err)
	}
	var reasoned bool // a '#' line is in scope until the next blank line
	var n int
	for i, raw := range strings.Split(string(b), "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case line == "":
			reasoned = false
		case strings.HasPrefix(line, "#"):
			// A divider such as "# --- Foo ---" is still prose a reader can use.
			reasoned = true
		default:
			n++
			if !reasoned {
				t.Errorf("allowed-drift.txt:%d: entry %q has no reason comment above it; "+
					"state which of the two tests it fails, with the query or file:line that shows it", i+1, line)
			}
		}
	}
	if n == 0 {
		t.Fatal("parsed no entries from allowed-drift.txt -- the parser or the file shape changed")
	}
	t.Logf("%d allowed-drift entries, all with reasons", n)
}
