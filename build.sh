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
  build_petrock
  cd tmp
  [[ -d blog ]] && { yes | rm -rf ./blog; }
  petrock new blog github.com/dhamidi/blog
  cd blog
  petrock feature posts
  go run ./cmd/blog serve --log-level=debug &
  sleep 1
  test_project_create_post
  wait %1
}

test_project_create_post() {
  curl 'http://localhost:8080/posts/new' \
    -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' \
    -H 'Accept-Language: en-US,en;q=0.9' \
    -H 'Cache-Control: max-age=0' \
    -H 'Connection: keep-alive' \
    -H 'Content-Type: application/x-www-form-urlencoded' \
    -H 'DNT: 1' \
    -H 'Origin: http://localhost:8080' \
    -H 'Referer: http://localhost:8080/posts/new' \
    -H 'Sec-Fetch-Dest: document' \
    -H 'Sec-Fetch-Mode: navigate' \
    -H 'Sec-Fetch-Site: same-origin' \
    -H 'Sec-Fetch-User: ?1' \
    -H 'Upgrade-Insecure-Requests: 1' \
    -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36' \
    -H 'sec-ch-ua: "Google Chrome";v="135", "Not-A.Brand";v="8", "Chromium";v="135"' \
    -H 'sec-ch-ua-mobile: ?0' \
    -H 'sec-ch-ua-platform: "macOS"' \
    --data-raw 'csrf_token=token&name=Test&description=Test'
}

main "$@"
