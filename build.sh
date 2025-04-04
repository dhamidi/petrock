#!/usr/bin/env bash

main() {
  step build_skeleton
  step build_petrock
  step test_petrock
}

step() {
  local name="$1"
  printf "START %s\n" "$name"
  if ! $name; then
    printf "FAIL %s\n" "$name"
    exit 1
  else
    printf "OK   %s\n" "$name"
  fi
}

build_skeleton() {
  cp internal/skeleton/go.mod{.skel,} && go build ./internal/...
}

build_petrock() {
  rm internal/skeleton/go.mod # to allow go:embed to do its work
  go build ./cmd/...
}

test_petrock() {
  ./petrock test
}

main "$@"
