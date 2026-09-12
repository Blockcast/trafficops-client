# Upstream Sync State

- Upstream repository: https://github.com/Blockcast/trafficcontrol
- Last-synced commit: unknown
- Last-synced date: unknown; the initial extraction was committed on 2025-12-07

The initial extraction did not record its Traffic Control source commit, and the
existing history does not establish one unambiguously. Do not substitute a
commit inferred from timestamps. The first selective upstream sync, expected to
produce `v0.1.6`, must establish the base by replacing the unknown values above
with the exact Traffic Control commit SHA and sync date.

Every later sync pull request must update this file. Once a base is established,
start the next review with:

```sh
git diff <last-synced-commit>..master -- lib/go-tc/
```

## Struct drift check

`tools/structdrift` compares struct *fields* between upstream `lib/go-tc` and
this repository's `go-tc`, reporting fields that exist upstream on a struct that
exists on both sides. It exists because a missing field is silent: a caller doing
a full-object read-modify-write discards it on decode and omits it on encode, and
Traffic Ops then binds SQL NULL over a column nobody meant to touch
(BLO-33496 erased `hardware_id` and `network_id` this way).

Run it locally against an upstream checkout:

```sh
go run ./tools/structdrift -upstream /path/to/trafficcontrol/lib/go-tc -vendored ./go-tc
```

Known omissions are baselined in `tools/structdrift/allowed-drift.txt`. That
baseline was recorded against trafficcontrol master
`76595609e1bce75234b0ce78557fb0654f60843d` on 2026-09-12. **This is not the
last-synced base** and must not be substituted for the unknown value above — no
selective sync was performed; only the three `ServerV40`/`ServerV50` fields named
in BLO-33496 were vendored.

Note that a whole-file diff is not a usable gate here: 38 of the 88 files present
in both trees differ, because this repository is a selective extraction.
