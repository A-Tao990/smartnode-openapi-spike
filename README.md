# smartnode-openapi-spike

A small, self-contained experiment: describe **one** read-only endpoint of the
[Rocket Pool smartnode](https://github.com/rocket-pool/smartnode) HTTP API
(`GET /api/network/node-fee`) in OpenAPI 3.0, generate Go types / server interface / client with
[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) v2.5.0, and check in CI that the generated code is
reproducible and matches a stubbed server's behaviour.

**This is not an official Rocket Pool project and has no endorsement from the smartnode maintainers.**

## Reproduce in 3 steps

Requires Go (the version in `go.mod`) and network access to fetch modules.

```sh
git clone https://github.com/A-Tao990/smartnode-openapi-spike && cd smartnode-openapi-spike
bash ./regen.sh --check        # regenerates into a temp file, compares with api/api.gen.go, modifies nothing
go vet ./... && go test -count=1 ./...
```

`bash ./regen.sh` (without `--check`) rewrites `api/api.gen.go`. Set `GO_BIN=/path/to/go` to use a specific Go binary
instead of the one on `PATH`. The generator is pinned through `tools.go` + `go.mod`.

## Contents

| Path | What |
|---|---|
| `openapi.yaml`, `oapi-config.yaml` | Hand-written contract for one endpoint, and generator config |
| `api/api.gen.go` | Generated code (do not edit) |
| `stub_test.go` | Contract tests against a **stub** server (200, 401 + `WWW-Authenticate`, non-GET, 500 error envelope) |
| `regen.sh` | Regenerate / `--check` reproducibility |
| `testdata/upstream/` | Verbatim upstream snapshots used as reference (not compiled); see `SOURCES.md` |
| `.github/workflows/ci.yml` | `regen.sh --check`, `go vet`, `go test` |

## Limitations

- **Single endpoint, extracted stub.** Only `/api/network/node-fee` is covered. The server in the tests is a hand-written
  stub that mimics the observable HTTP behaviour (route, bearer-token check, JSON envelope) of the upstream code.
- **The real `snroute` was not compiled and smartnode was not migrated.** Nothing here builds or runs against the real
  daemon, go-ethereum, or a chain. The upstream files in `testdata/upstream/` are only read, never built.
- **Fake data.** The fee values returned by the stub are hard-coded.
- **Method mismatch.** The real router registers handlers without restricting the HTTP method, while the generated
  standard-library mux binds the operation to `GET`. Adopting the generated router would tighten behaviour.
- **Auth and scopes are not integrated.** Upstream's token handling (including the `unauthenticatedReads` setting and
  write scopes) is not implemented by generated code; the test's `requireToken` middleware covers only a subset. The
  `security` section documents upstream semantics; code generation does not implement authentication.
- Other endpoints (including ones using `*big.Int` types and write routes) are not modelled.
- "Tests pass" means the generated code agrees with the stub, not with the real smartnode.

## License

GPL-3.0 — see `LICENSE`. It is the license of the upstream project; `openapi.yaml` is derived from upstream behaviour,
`api/api.gen.go` is generated from it, and the new example code in this repository is also released under GPL-3.0.
Upstream sources, commit and checksums are listed in `SOURCES.md`.

## AI assistance disclosure

This repository was prepared with the help of AI coding assistants (Claude and Codex), including the OpenAPI
contract, the stub tests and the scripts. Output was run locally and reviewed by the repository owner's team, but it
may contain mistakes.
