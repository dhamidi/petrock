# Overview

Read through docs/high-level.md and internal/skeleton/feature_template.

We'll be mostly working with these files.

The goal for this project is to make the default feature template more enticing.

Ultimately the following should be possible after running `petrock feature posts` (`posts` corresponds to `petrock_example_feature_name` in internal/skeleton/feature_template)

These routes will be defined:

- `GET /posts/new` renders a form using (core/form.go) for adding a new post
- `POST /posts/new` accepts the form, converts it into a `CreateCommand`, and dispatches it to the core.Executor
  - during validation, this command trims all incoming strings and makes sure they are not empty
  - it also checks that the provided post string ID is unique
  - in case of errors the corresponding http handler will render the form with errors
  - in case of success, it redirects to `GET /posts` with HTTP Status SeeOther
- `GET /posts` lists all posts by issuing a posts.ListQuery
  - it includes a button which links to the new posts form
- `GET /posts/{id}` renders a given post as a series of labels + pre-formatted text fields for each attribute on a post
- `GET /posts/{id}/edit` renders a form for editing the `content` field of a post
- `GET /posts/{id}/delete` renders a form with a button that will trigger a `DeleteCommand` for the given post
- `POST /posts/{id}/delete` will accept the form from the previous route and actually dispatch the `DeleteCommand`
