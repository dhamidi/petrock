#!/usr/bin/env bash

# Helper function to print error messages and exit
error_exit() {
  printf "Error: %s\n" "$1" >&2
  exit 1
}

main() {
  local log_level="info"
  local target_step=""
  local remaining_args=()

  # Parse arguments
  for arg in "$@"; do
    case "$arg" in
    --debug)
      log_level="debug"
      ;;
    *)
      # Collect non-flag arguments
      remaining_args+=("$arg")
      ;;
    esac
  done

  # Determine the target step, if any
  if [[ ${#remaining_args[@]} -gt 1 ]]; then
    error_exit "Too many steps specified: ${remaining_args[*]}"
  elif [[ ${#remaining_args[@]} -eq 1 ]]; then
    target_step="${remaining_args[0]}"
  fi

  # Set log level environment variable
  export PETROCK_LOG_LEVEL="$log_level"

  # Execute steps
  if [[ -n "$target_step" ]]; then
    # Check if the target step function exists
    if declare -F "$target_step" >/dev/null; then
      step "$target_step"
    else
      error_exit "Unknown step: '$target_step'"
    fi
  else
    # Run all default steps
    step build_skeleton
    step build_petrock
    step test_petrock
  fi
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
  rm -f internal/skeleton/go.mod # to allow go:embed to do its work
  go build ./cmd/...
  go install ./cmd/petrock
}

test_petrock() {
  ./petrock test
}

test_project() {
  cd tmp
  [[ -d blog ]] && { yes | rm -rf ./blog; }
  petrock new blog github.com/dhamidi/blog
  cd blog
  petrock feature posts
  go run ./cmd/blog serve
}

main "$@"
