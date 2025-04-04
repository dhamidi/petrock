#!/usr/bin/env bash

main() {
  build_skeleton
  build_petrock
  test_petrock
}

build_skeleton() {
  cp internal/skeleton/go.mod{.skel,} && go build internal/skeleton
}

build_petrock() {
  rm internal/skeleton/go.mod # to allow go:embed to do its work
  go build ./cmd/...
}

test_petrock() {
  ./petrock test
}

main "$@"
