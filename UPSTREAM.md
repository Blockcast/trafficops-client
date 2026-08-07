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
