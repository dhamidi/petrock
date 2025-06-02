We need to refactor how I/O is done with the user in ./cmd/ and ./internal/skeleton/cmd/petrock_example_project_name/

What I want:

- a UI interface with methods to present information to the user, and prompt the user for information,
- all information that is relevant to the user should be presented through using this interface,
- commands invoked by the user should never _log_ output, they should communicate with the user only through the UI interface.

Please make a plan in docs/plans/text-ui/plan.md for making the necessary changes.

In phase 1, we should only apply the changes to ./cmd

In phase 2, we'll rework the command line tools in ./internal/skeleton/cmd/petrock_example_project_name
