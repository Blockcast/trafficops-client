# Traffic Ops Client

This module is Blockcast's import-independent fork of the Apache Traffic Control
Go client libraries. It is consumed by Magma and is released separately from
Traffic Control.

## Ownership

The **GStack Release Engineer** owns releases of this repository. The release
owner is responsible for reviewing upstream sync scope, merging a green sync,
tagging the release, and opening the corresponding Magma dependency bump.

Upstream changes are not published here automatically. A change to Traffic
Control's `lib/go-tc` does not reach Magma until this repository is reconciled,
tagged, and Magma is updated.

## Upstream Re-sync

The fork is intentionally independent, but it is not intentionally frozen.
Re-syncs are selective reconciliations, not directory replacements.

1. Start from the latest `main` and record the Traffic Control source commit in
   the pull request and sync commit. Update `UPSTREAM.md` with that exact commit
   and the sync date.
2. Inventory the overlap between upstream `lib/go-tc/` and this repository's
   `go-tc/`. Every upstream file whose relative path already exists in `go-tc/`
   is in scope for review. New upstream-only files are added only when a
   supported Blockcast consumer needs them.
3. Review each upstream diff and preserve intentional fork changes, including
   the `github.com/Blockcast/trafficops-client` module path, Blockcast API
   extensions, and Blockcast dependency substitutions. Rewrite imports from
   Traffic Control packages to the corresponding package in this module; do
   not restore `github.com/apache/trafficcontrol` imports.
4. Include related changes from the other extracted libraries (`go-util/`,
   `go-rfc/`, `go-log/`, `go-llog/`, `toclientlib/`, `v4-client/`, and
   `v5-client/`) when the selected `go-tc` changes depend on them.
5. Run `go build ./...` and `go test ./...`. Merge only after the repository CI
   workflow passes and the reconciliation has been reviewed.

The pull request must list files intentionally omitted from the sync and why.
That list distinguishes a deliberate fork decision from accidental drift.

As verified on 2026-08-07, this fork's `go-tc/` is a strict subset of upstream
`lib/go-tc/`: 165 files overlap, 112 differ, 53 are byte-identical, 38 exist
only upstream, and 0 exist only in the fork. The fork has no Blockcast-original
`go-tc` files; its divergence is within files, primarily from import rewrites.
Maintain the 0-fork-only invariant. Every sync pull request must report the file
set comparison and explicitly explain any change to that invariant.

## Deletion Pre-check

Before deleting a symbol from Traffic Control's `lib/go-tc/`, the engineer who
owns the deletion must run both cross-repository checks:

1. Search this repository's `go-tc/` for the symbol and its serialized values.
2. Search Blockcast Magma's `cdn/` tree for the same identifiers and values.

Record both results in the Traffic Control deletion pull request. A local search
in Traffic Control cannot clear either consumer. The GStack Release Engineer
reviews a positive fork match to determine whether a selective sync, release,
and Magma dependency bump are required before the deletion can land.

Check absence as well as presence. For example, `MulticastServerType` and
`AMT_RELAY` occur in neither this fork nor Magma's `cdn/`, so that deletion has
no Magma exposure through this module. Do not assume that result for other
symbols, and do not assume a fork match implies a Magma use without the second
search.

## Release and Magma Bump

After a re-sync is merged and CI is green, the GStack Release Engineer:

1. Chooses the next semantic version and creates an annotated tag on the green
   `main` commit, for example `git tag -a v0.1.6 -m "Release v0.1.6"`.
2. Pushes the tag with `git push origin v0.1.6` and verifies that the tag points
   to the tested commit.
3. Opens a Magma pull request updating every `trafficops-client` version in
   `cdn/cloud/go/go.mod`, including both direct requirements and replacement
   directives, to the same new tag.
4. Runs the Magma CDN Go build and tests before merging the dependency bump.

Magma must not be bumped before the tag exists and this repository's CI is
green. The Traffic Control change is not delivered to Magma until that bump is
merged.
