# Shell Script Style Guide

These rules apply to all shell scripts (`*.sh`) within this project, particularly those intended for automation, build processes, or command-line tooling. The goal is to ensure scripts are readable, maintainable, portable, and robust.

# Rule 1: Use `#!/usr/bin/env bash`

Start scripts with `#!/usr/bin/env bash` for portability, ensuring the script uses the `bash` interpreter found in the user's `PATH`.

```bash
#!/usr/bin/env bash
# Good: Uses env to find bash in the user's PATH

echo "Hello"
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
  echo "Processing arguments: $@"
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
  echo "Building..."
  # build commands
}

deploy_project() {
  echo "Deploying..."
  # deploy commands
}

main() {
  build_project
  deploy_project
}

# Good: Logic is broken down into functions.
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
  echo "$my_var"
}

main() {
  local my_var="main_var"
  my_func
  echo "$my_var" # Output: main_var (not overwritten by my_func)
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
  # Good: Quoting prevents word splitting.
  touch "$filename"
  ls -l "$filename"

  # Good: "$@" expands each argument as a separate word.
  for arg in "$@"; do
    echo "Arg: $arg"
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
  echo "Attempting operation..."
  if ! mkdir /nonexistent_dir/subdir; then
    # Good: Checks exit status, reports error to stderr, exits non-zero.
    echo "Error: Failed to create directory." >&2
    exit 1
  fi
  echo "Operation successful."
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
  echo "Building..."
}

cmd_test() {
  echo "Testing..."
}

main() {
  local subcommand="$1"
  shift # Remove subcommand from argument list

  # Good: Clear dispatch based on the first argument.
  case "$subcommand" in
    build)
      cmd_build "$@"
      ;;
    test)
      cmd_test "$@"
      ;;
    *)
      echo "Error: Unknown subcommand '$subcommand'" >&2
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
  # Good: Checks if 'goimports' exists before trying to use it.
  if ! command -v goimports &> /dev/null; then
    echo "Error: goimports is not installed. Please install it." >&2
    exit 1
  fi

  echo "Running goimports..."
  # goimports ...
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
