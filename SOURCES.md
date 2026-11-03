# Upstream sources

Snapshots in `testdata/upstream/` are verbatim copies of files from
[rocket-pool/smartnode](https://github.com/rocket-pool/smartnode) (GPL-3.0), pinned to commit
`72b831c44b5b93fcc0e16fadabff610592406197` (master, committed 2026-11-03T02:52:14Z). Fetched on **2026-11-03**.
They are reference material only: they live under `testdata/`, which the Go toolchain ignores, so they are never compiled here.
The upstream copies carry no per-file license headers; the license is the repository's `LICENSE` (copied as the root `LICENSE` of this repo).

| Local file (testdata/upstream/) | Upstream path | SHA-256 |
|---|---|---|
| snroute.go | rocketpool/api/snroute/snroute.go | 10c9c2319e28308ba04aed9cd8b64847e910642f0e7459b184c4c133ce42c1be |
| network_routes.go | rocketpool/api/network/routes.go | 7831c0e4d5e8cba8385404b896f2fbed6078152b42c00472e56aabfb7e327c7f |
| node-fee.go | rocketpool/api/network/node-fee.go | 9f9cd1b65f5c3bd4fec35e3e32d8c57efb32f514b2f846e0df5cafdfc45db2d9 |
| response.go | rocketpool/api/response/response.go | 5ea2ea161c24eba23935ea014e1a73db47387e0107fbc39fb0d78e125e8653ba |
| api_network.go | shared/types/api/network.go | ea261126d44a88a18d2adf17352388c2d2e3aaa72ad3f4e9a5568c3eae8b5be0 |
| LICENSE | LICENSE | 8ceb4b9ee5adedde47b31e975c1d90c73ad27b6b165a1dcd80c7c545eb65b903 |

Original URL pattern: `https://raw.githubusercontent.com/rocket-pool/smartnode/72b831c44b5b93fcc0e16fadabff610592406197/<upstream path>`

Verify: compare `sha256sum testdata/upstream/*` with the table above.

Derived files in this repo: `openapi.yaml` (hand-written from the behaviour of the files above) and `api/api.gen.go` (generated from it).

Related upstream client (not copied here): [shared/services/rocketpool/network.go](https://github.com/rocket-pool/smartnode/blob/72b831c44b5b93fcc0e16fadabff610592406197/shared/services/rocketpool/network.go) calls the node-fee endpoint with GET.
