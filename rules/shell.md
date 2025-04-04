# Shell Script Style Guide

These rules apply to all shell scripts (`*.sh`) within this project, particularly those intended for automation, build processes, or command-line tooling. The goal is to ensure scripts are readable, maintainable, portable, and robust.

# Rule 1: Use `#!/usr/bin/env bash`

Start scripts with `#!/usr/bin/env bash` for portability, ensuring the script uses the `bash` interpreter found in the user's `PATH`.

```bash
#!/usr/bin/env bash
# Good: Uses env to find bash, has main, uses printf.

main() {
  printf "Hello\n"
}

main "$@"
```

```bash
#!/bin/bash
# Bad: Hardcodes the path to bash, which might not exist at /bin/bash

echo "Hello"
```

# Rule 2: Use a `main()` function

Use a `main()` function as the primary entry point and call it with `main "$@"` at the script's end to pass all arguments correctly. This improves structure and readability.

```bash
#!/usr/bin/env bash

main() {
  # Good: Uses local, quotes arguments, uses printf.
  local arg
  printf "Processing arguments:\n"
  for arg in "$@"; do
    printf " - %s\n" "$arg"
  done
}

# Good: Script logic is contained in main, arguments are passed correctly.
main "$@"
```

```bash
#!/usr/bin/env bash

# Bad: Script logic is at the top level, less organized.
# Argument handling might be less explicit.
echo "Processing arguments: $@"

# No main function call.
```

# Rule 3: Decompose logic into functions

Decompose logic into well-named functions (e.g., `verb_noun` or using prefixes like `cmd_` for subcommands). This improves modularity and makes the script easier to understand and test.

```bash
#!/usr/bin/env bash

build_project() {
  # Good: Uses printf, includes placeholder error check.
  printf "Building...\n"
  # build commands
  # if ! some_build_command; then
  #   printf "Error: Build failed.\n" >&2
  #   return 1 # Indicate failure to the caller
  # fi
  return 0 # Indicate success
}

deploy_project() {
  # Good: Uses printf, includes placeholder error check.
  printf "Deploying...\n"
  # deploy commands
  # if ! some_deploy_command; then
  #   printf "Error: Deploy failed.\n" >&2
  #   return 1 # Indicate failure to the caller
  # fi
  return 0 # Indicate success
}

main() {
  # Good: Checks return status of functions.
  if ! build_project; then
    exit 1
  fi
  if ! deploy_project; then
    exit 1
  fi
}

# Good: Logic is broken down into functions, main called correctly.
main "$@"
```

```bash
#!/usr/bin/env bash

main() {
  # Bad: All logic is crammed into one function or the top level.
  echo "Building..."
  # build commands

  echo "Deploying..."
  # deploy commands
}

main "$@"
```

# Rule 4: Use `local` for variables

Declare variables within functions using `local` to limit their scope. This prevents accidental modification of global variables or variables in calling functions.

```bash
#!/usr/bin/env bash

my_func() {
  local my_var="hello" # Good: Variable scope is limited to my_func.
  # Good: Uses printf.
  printf "my_func variable: %s\n" "$my_var"
}

main() {
  local my_var="main_var"
  my_func
  # Good: Uses printf.
  printf "main variable: %s\n" "$my_var" # Output: main variable: main_var
}

main "$@"
```

```bash
#!/usr/bin/env bash

my_func() {
  my_var="hello" # Bad: Variable is global by default.
  echo "$my_var"
}

main() {
  my_var="main_var"
  my_func
  echo "$my_var" # Output: hello (overwritten by my_func)
}

main "$@"
```

# Rule 5: Quote variable expansions

Always quote variable expansions (e.g., `"$variable"`, `"$@"`) to prevent unexpected word splitting and filename generation (globbing) based on whitespace or special characters in the variable's value.

```bash
#!/usr/bin/env bash

main() {
  local filename="file with spaces.txt"
  local arg

  # Good: Quoting prevents word splitting.
  # Good: Includes error checking.
  if ! touch "$filename"; then
    printf "Error: Failed to touch '%s'.\n" "$filename" >&2
    # Decide whether to exit or continue
  else
    # Good: Includes error checking.
    if ! ls -l "$filename"; then
       printf "Error: Failed to list '%s'.\n" "$filename" >&2
       # Decide whether to exit or continue
    fi
  fi


  # Good: "$@" expands each argument as a separate word.
  # Good: Uses printf.
  printf "\nArguments:\n"
  for arg in "$@"; do
    printf " - %s\n" "$arg"
  done
}

main "arg 1" "arg 2"
```

```bash
#!/usr/bin/env bash

main() {
  local filename="file with spaces.txt"
  # Bad: Unquoted variable undergoes word splitting.
  # This tries to touch "file", "with", and "spaces.txt".
  # touch $filename # This would likely fail or create wrong files.

  # Bad: $* expands all arguments into a single word.
  # Bad: $@ (unquoted) splits arguments containing spaces.
  for arg in $@; do
    echo "Arg: $arg" # Output might be split unexpectedly if args have spaces.
  done
}

main "arg 1" "arg 2"
```

# Rule 6: Implement robust error handling

Check command exit statuses (`$?` or `if ! command`), exit with a non-zero status on failure (`exit 1`), and print error messages to standard error (`>&2`). Use `set -e` cautiously, as it can sometimes mask errors or make debugging harder. Explicit checks are often clearer.

```bash
#!/usr/bin/env bash

main() {
  # Good: Uses printf.
  printf "Attempting operation...\n"
  # Example using a directory likely to be writable
  local target_dir="./temp_test_dir"

  # Good: Checks exit status, reports error using printf to stderr, exits non-zero.
  if ! mkdir "$target_dir"; then
    printf "Error: Failed to create directory '%s'.\n" "$target_dir" >&2
    exit 1
  fi
  printf "Operation successful. Directory '%s' created.\n" "$target_dir"

  # Clean up the created directory (optional, for example purposes)
  rmdir "$target_dir"
}

main "$@"
```

```bash
#!/usr/bin/env bash

main() {
  echo "Attempting operation..."
  # Bad: No check for command success. Script might continue after failure.
  mkdir /nonexistent_dir/subdir
  # Bad: Error message (if any) goes to stdout.
  echo "Failed to create directory." # Incorrectly implies success check happened.
  # Bad: Script exits 0 even on failure.
  echo "Operation finished."
}

main "$@"
```

# Rule 7: Prefer `printf` over `echo`

Prefer `printf` over `echo` for more consistent and controllable output formatting, especially when dealing with variables that might contain special characters or backslashes. `echo` behavior can vary between shells and versions.

```bash
#!/usr/bin/env bash

main() {
  local name="User"
  local message="Welcome!\nCheck your settings."

  # Good: printf handles formatting and special characters predictably.
  printf "Hello, %s.\n" "$name"
  printf -- "---\n%s\n---\n" "$message" # Handles newline correctly
}

main "$@"
```

```bash
#!/usr/bin/env bash

main() {
  local name="User"
  local message="Welcome!\nCheck your settings."

  # Bad: echo's handling of options (-n, -e) and backslashes varies.
  echo "Hello, $name."
  echo "$message" # May print \n literally depending on echo version/flags.
}

main "$@"
```

# Rule 8: Use parameter expansion for defaults

Use parameter expansion for setting default values (e.g., `local var="${1:-default_value}"`). This is concise and standard.

```bash
#!/usr/bin/env bash

main() {
  # Good: Sets log_level to $1 if provided, otherwise defaults to 'info'.
  local log_level="${1:-info}"
  printf "Log level: %s\n" "$log_level"
}

main "$@"       # Output: Log level: info
main "debug"    # Output: Log level: debug
```

```bash
#!/usr/bin/env bash

main() {
  local log_level
  if [ -n "$1" ]; then
    log_level="$1"
  else
    log_level="info"
  fi
  # Bad: More verbose than parameter expansion.
  printf "Log level: %s\n" "$log_level"
}

main "$@"
main "debug"
```

# Rule 9: Structure subcommands clearly

For command-line tools, structure subcommands clearly, often using a `case` statement in `main` or dedicated `cmd_subcommand` functions. This makes the script's interface understandable and extensible.

```bash
#!/usr/bin/env bash

cmd_build() {
  # Good: Uses printf. Placeholder for actual build logic.
  printf "Building...\n"
  # Add build commands and error checking here
  return 0 # Indicate success
}

cmd_test() {
  # Good: Uses printf. Placeholder for actual test logic.
  printf "Testing...\n"
  # Add test commands and error checking here
  return 0 # Indicate success
}

main() {
  local subcommand="${1:-}" # Use parameter expansion for safety
  shift || true # Shift arguments, ignore error if no arguments

  # Good: Clear dispatch based on the first argument.
  case "$subcommand" in
    build)
      # Good: Calls function and checks status.
      if ! cmd_build "$@"; then exit 1; fi
      ;;
    test)
      # Good: Calls function and checks status.
      if ! cmd_test "$@"; then exit 1; fi
      ;;
    "")
      # Handle empty subcommand case if necessary
      printf "Error: No subcommand provided.\n" >&2
      # Potentially show help here
      exit 1
      ;;
    *)
      # Good: Uses printf for error message.
      printf "Error: Unknown subcommand '%s'\n" "$subcommand" >&2
      exit 1
      ;;
  esac
}

main "$@"
```

```bash
#!/usr/bin/env bash

main() {
  # Bad: Logic based on flags or argument position is less clear for distinct actions.
  if [[ "$1" == "--build" ]]; then
    echo "Building..."
  elif [[ "$1" == "--test" ]]; then
    echo "Testing..."
  else
    echo "Error: Specify --build or --test" >&2
    exit 1
  fi
}

main "$@"
```

# Rule 10: Check for required commands

Check for the existence of required external commands using `command -v command_name &> /dev/null` before attempting to use them. This provides clearer error messages to the user if dependencies are missing.

```bash
#!/usr/bin/env bash

main() {
  local required_cmd="goimports"
  # Good: Checks if command exists before trying to use it.
  # Good: Uses printf for error message.
  if ! command -v "$required_cmd" &> /dev/null; then
    printf "Error: Required command '%s' is not installed. Please install it.\n" "$required_cmd" >&2
    exit 1
  fi

  # Good: Uses printf for status message.
  printf "Running %s...\n" "$required_cmd"
  # "$required_cmd" ... # Example usage
}

main "$@"
```

```bash
#!/usr/bin/env bash

main() {
  # Bad: Assumes 'goimports' exists. If not, the script fails with a potentially
  # confusing "command not found" error later on.
  echo "Running goimports..."
  goimports ... # This line will fail if goimports isn't installed.
}

main "$@"
```
