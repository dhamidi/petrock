#!/usr/bin/env bash

# Function to lint Go files
lint_go() {
    local file="$1"
    # echo "Linting Go file: $file" # Removed for less verbosity
    goimports -w "$file"
    # Add other Go linting commands here if needed
    # e.g., golangci-lint run "$file"
    return $? # Return the exit code of the last command
}

main() {
    local exit_code=0

    # Check if goimports is installed
    if ! command -v goimports &> /dev/null; then
        echo "Error: goimports is not installed. Please install it (go install golang.org/x/tools/cmd/goimports@latest)" >&2
        exit 1
    fi

    for file in "$@"; do
        # echo "Processing file: $file" # Removed for less verbosity
        local extension="${file##*.}"

        case "$extension" in
            go)
                lint_go "$file"
                ;;
            # Add cases for other file types here
            # sh)
            #    lint_sh "$file"
            #    ;;
            *)
                # echo "Skipping file with unknown or unhandled extension: $file" # Removed for less verbosity
                # No action needed, implicitly exits with 0 for this file
                ;;
        esac

        # Capture the exit code of the linter function
        local current_exit_code=$?
        if [ $current_exit_code -ne 0 ]; then
            echo "Error linting file: $file (Exit code: $current_exit_code)" >&2
            exit_code=1 # Set the overall exit code to non-zero if any linter fails
        fi
    done

    exit $exit_code
}

main "$@"
