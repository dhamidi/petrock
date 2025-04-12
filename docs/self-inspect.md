# Self Inspection

Petrock applications include a `self inspect` command that provides detailed information about the application structure, including:

- All registered commands with their JSON schema
- All registered queries with their JSON schema
- All HTTP routes
- All features

## Usage

```shell
# Basic usage (outputs JSON to stdout)
$ myapp self inspect

# Specify database path
$ myapp self inspect --db-path=custom.db
```

## Output Format

The command outputs structured JSON that describes the application. Here's an example snippet:

```json
{
  "commands": [
    {
      "name": "posts/create",
      "type": "posts.CreatePostCommand",
      "properties": {
        "title": { "type": "string", "description": "Title of the post" },
        "content": { "type": "string", "description": "Content of the post" },
        "authorID": { "type": "string", "description": "ID of the post author" }
      },
      "required": ["title", "content", "authorID"]
    }
  ],
  "queries": [
    {
      "name": "posts/list",
      "type": "posts.ListPostsQuery",
      "properties": {
        "page": { "type": "integer", "default": 1 },
        "pageSize": { "type": "integer", "default": 10 },
        "authorIDFilter": { "type": "string" }
      },
      "required": ["page", "pageSize"],
      "result": {
        "type": "object",
        "properties": {
          "posts": { "type": "array" },
          "totalCount": { "type": "integer" },
          "page": { "type": "integer" },
          "pageSize": { "type": "integer" }
        }
      }
    }
  ],
  "routes": [
    "GET /",
    "GET /commands",
    "POST /commands",
    "GET /queries",
    "GET /queries/{feature}/{queryName}",
    "GET /posts",
    "GET /posts/{id}"
  ],
  "features": [
    "posts"
  ]
}
```

## Implementation Details

The self-inspection feature is implemented using the following components:

1. **App Structure**: The core `App` struct tracks features, routes, and dependencies centrally
2. **Command Registry**: Provides information about registered commands
3. **Query Registry**: Provides information about registered queries
4. **Route Tracking**: Routes are tracked when registered via `app.RegisterRoute`
5. **Feature Tracking**: Features are tracked when registered via `app.RegisterFeature`

### Command Schema Generation

Command schemas are generated automatically using reflection. For each registered command:

1. Command type is inspected to extract field information
2. Field types are converted to JSON Schema types
3. Properties are extracted from struct fields
4. Required fields are determined

### Query Schema Generation

Query schemas include both input parameters and result schema:

1. Query input parameters are extracted from struct fields
2. Result type is included to describe the expected response format
3. Field types are converted to appropriate JSON Schema types

## Use Cases

- **API Discovery**: Understand available commands and queries
- **Client Generation**: Auto-generate client code for APIs
- **Documentation**: Generate up-to-date API documentation
- **Testing**: Generate test cases for commands and queries
- **Validation**: Verify application structure and dependencies