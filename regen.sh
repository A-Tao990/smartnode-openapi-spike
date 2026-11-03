#!/usr/bin/env bash
# Regenerate api/api.gen.go from openapi.yaml with oapi-codegen v2.5.0 (pinned in go.mod via tools.go).
#
# Usage:
#   ./regen.sh          regenerate api/api.gen.go in place
#   ./regen.sh --check  generate into a temp file and compare; checked-in files are NOT modified.
#                       Exit 0 = identical, 1 = differs or api/api.gen.go is missing.
#
# Go toolchain: uses `go` from PATH, or set GO_BIN=/path/to/go to override.
# GOPATH/GOCACHE/PATH are left as the caller has them.
set -euo pipefail
cd "$(dirname "$0")"

mode="${1:-}"
case "$mode" in
  ""|--check) ;;
  *) echo "usage: $0 [--check]" >&2; exit 2 ;;
esac

GO="${GO_BIN:-go}"
if ! command -v "$GO" >/dev/null 2>&1; then
  echo "error: Go toolchain not found ('$GO'). Install Go or set GO_BIN=/path/to/go" >&2
  exit 2
fi

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

# Generate to a temp output; never touch the checked-in file in --check mode.
# The output path is passed with -o (oapi-config.yaml deliberately has no `output:` key).
"$GO" run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
  -config oapi-config.yaml -o "$tmp_dir/api.gen.go" openapi.yaml

if [ "$mode" = "--check" ]; then
  if [ ! -f api/api.gen.go ]; then
    echo "FAIL: api/api.gen.go is missing; run ./regen.sh to create it" >&2
    exit 1
  fi
  if diff -u api/api.gen.go "$tmp_dir/api.gen.go"; then
    echo "OK: api/api.gen.go is up to date"
  else
    echo "FAIL: api/api.gen.go differs from openapi.yaml; run ./regen.sh" >&2
    exit 1
  fi
else
  mkdir -p api
  cp "$tmp_dir/api.gen.go" api/api.gen.go
  echo "wrote api/api.gen.go"
fi
