// Command structdrift reports struct fields that exist in Traffic Control's
// lib/go-tc but are missing from this repository's vendored go-tc copy.
//
// This repository is a selective extraction of upstream, so the two trees
// legitimately differ in most files — a whole-file diff produces ~38 differing
// files out of 88 and is unusable as a gate. What is never legitimate is a
// struct that exists on both sides carrying a field on only one of them: a
// caller doing a full-object read-modify-write through this client silently
// discards the field on decode, omits it on encode, and Traffic Ops binds SQL
// NULL over a column nobody meant to touch. That is BLO-33496, which erased the
// orc8r gateway-identity columns hardware_id and network_id.
//
// Usage:
//
//	structdrift -upstream /path/to/trafficcontrol/lib/go-tc -vendored ./go-tc
//
// WHERE THIS RUNS. Not here. Blockcast/trafficcontrol is private and this
// repository is public, so a job here cannot read upstream without a
// Blockcast/trafficcontrol credential living in a public repository's Actions
// secrets. The check runs instead from trafficcontrol's own
// .github/workflows/go-tc-drift.yml, which has lib/go-tc natively and checks
// this repository out for the vendored tree and for this tool -- reading a
// public repository needs no secret, so that direction costs no credential at
// all. See BLO-33506.
//
// At PR time here, go-tc/roundtrip_drift_test.go covers the fields that were
// measured to be erasing; it is a plain go test and needs no upstream access.
//
// Exits non-zero when drift is found. Known-and-accepted omissions are listed
// in allowed-drift.txt as "StructName.FieldName" lines, each with a reason.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// fields returns every exported field of every top-level struct in dir, keyed
// by "StructName.FieldName". Embedded fields are keyed by their type name.
func fields(dir string) (map[string]bool, error) {
	out := map[string]bool{}
	paths, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no .go files in %s", dir)
	}
	fset := token.NewFileSet()
	for _, p := range paths {
		if strings.HasSuffix(p, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", p, err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok || !ts.Name.IsExported() {
				return true
			}
			for _, fld := range st.Fields.List {
				if len(fld.Names) == 0 { // embedded
					if id, ok := fld.Type.(*ast.Ident); ok {
						out[ts.Name.Name+"."+id.Name] = true
					}
					continue
				}
				for _, name := range fld.Names {
					if name.IsExported() {
						out[ts.Name.Name+"."+name.Name] = true
					}
				}
			}
			return true
		})
	}
	return out, nil
}

func allowed(path string) (map[string]bool, error) {
	out := map[string]bool{}
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out[line] = true
	}
	return out, s.Err()
}

func main() {
	upstream := flag.String("upstream", "", "path to trafficcontrol/lib/go-tc")
	vendored := flag.String("vendored", "./go-tc", "path to this repo's go-tc")
	allowFile := flag.String("allow", "tools/structdrift/allowed-drift.txt", "known-accepted omissions")
	flag.Parse()

	if *upstream == "" {
		fmt.Fprintln(os.Stderr, "-upstream is required")
		os.Exit(2)
	}

	up, err := fields(*upstream)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	mine, err := fields(*vendored)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	allow, err := allowed(*allowFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Only structs present on BOTH sides are compared. A struct this repo never
	// extracted is not drift, it is scope.
	vendoredStructs := map[string]bool{}
	for k := range mine {
		vendoredStructs[strings.SplitN(k, ".", 2)[0]] = true
	}

	var missing []string
	for k := range up {
		if mine[k] || allow[k] {
			continue
		}
		if vendoredStructs[strings.SplitN(k, ".", 2)[0]] {
			missing = append(missing, k)
		}
	}
	sort.Strings(missing)

	if len(missing) == 0 {
		fmt.Println("no struct drift")
		return
	}
	fmt.Fprintf(os.Stderr, "%d field(s) exist upstream but are missing from the vendored copy.\n", len(missing))
	fmt.Fprintf(os.Stderr, "A full-object read-modify-write through this client will silently NULL these columns.\n")
	fmt.Fprintf(os.Stderr, "Add the field, or record it in %s with a reason.\n\n", *allowFile)
	for _, m := range missing {
		fmt.Fprintf(os.Stderr, "  %s\n", m)
	}
	os.Exit(1)
}
