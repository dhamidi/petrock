# PetRock Project API Usage Guide

This guide explains how to interact with the JSON API exposed by the PetRock application server (typically started with `go run ./cmd/blog serve`).

## Mode of Operation: Commands and Queries

The API follows a Command Query Responsibility Segregation (CQRS) pattern:

1.  **Commands**: Used to *change* the state of the application (e.g., create, update, delete data).
    *   Sent via `POST` requests to the `/commands` endpoint.
    *   The specific command to execute is identified by a `type` field in the JSON body.
    *   Command-specific data is nested within a `payload` field in the JSON body.
    *   Commands are validated, logged, and then applied to the application's state.
    *   Successful command execution typically returns a `{"status":"success"}` JSON response with HTTP status `200 OK` or `202 Accepted`. Validation or processing errors usually result in `400 Bad Request` or `500 Internal Server Error`.

2.  **Queries**: Used to *read* the state of the application without changing it.
    *   Sent via `GET` requests to specific query endpoints.
    *   The endpoint path follows the pattern `/queries/{feature-name}/{query-name}` (e.g., `/queries/petrock_example_feature_name/list`).
    *   Query parameters (like IDs, filters, pagination) are passed as URL query string parameters (e.g., `?ID=some-id&page=2`).
    *   Successful queries return the requested data as a JSON response with HTTP status `200 OK`. Errors like "not found" might return `404 Not Found`, while other processing errors might return `500 Internal Server Error`.

## Important Data Structures

### Command Request Body (`POST /commands`)

```json
{
  "type": "feature-name/command-name",
  "payload": {
    "command_field_1": "value1",
    "command_field_2": 123
    // ... other command-specific fields
  }
}
```

*   `type`: (String, Required) The unique identifier for the command (e.g., `petrock_example_feature_name/create`).
*   `payload`: (Object, Required) An object containing the data required by the specific command. The fields within `payload` depend on the command definition.

### Query Request (GET `/queries/{feature-name}/{query-name}`)

Queries are invoked via GET requests, and their parameters are passed in the URL query string.

Example: `GET /queries/petrock_example_feature_name/get?ID=item-123`

The specific parameters accepted depend on the query definition (e.g., `ID`, `page`, `pageSize`, `filter`). Parameter names typically match the field names in the corresponding query struct in the code (case-sensitive, though the server might attempt case-insensitive matching for convenience).

## API Interaction Examples (`curl`)

*(Ensure the server is running, e.g., `go run ./cmd/blog serve`)*

### 1. List Available Commands

```bash
curl http://localhost:8080/commands
```

*(Expected Output: A JSON array of registered command names, e.g., `["petrock_example_feature_name/create", "petrock_example_feature_name/update", ...]`)*

### 2. List Available Queries

```bash
curl http://localhost:8080/queries
```

*(Expected Output: A JSON array of registered query names, e.g., `["petrock_example_feature_name/get", "petrock_example_feature_name/list", ...]`)*

### 3. Execute a Command (Create Example)

This example executes the `petrock_example_feature_name/create` command. Replace the `payload` fields with the actual fields required by your command.

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"type": "petrock_example_feature_name/create", "payload": {"name": "My First Item", "description": "Details about the item.", "created_by": "api_user"}}' \
  http://localhost:8080/commands
```

*(Expected Output on Success: `{"status":"success"}`)*

### 4. Execute a List Query

This example executes the `petrock_example_feature_name/list` query to retrieve a list of items.

```bash
curl http://localhost:8080/queries/petrock_example_feature_name/list
```

*(Expected Output: A JSON object containing an array of items and pagination details, e.g., `{"items":[...],"total_count":1,"page":1,"page_size":20}`)*

You can add query parameters for pagination or filtering if the query supports them:

```bash
# Example: Get page 2 with 5 items per page
curl "http://localhost:8080/queries/petrock_example_feature_name/list?page=2&pageSize=5"

# Example: Filter (if supported by the query handler)
curl "http://localhost:8080/queries/petrock_example_feature_name/list?filter=keyword"
```

### 5. Execute a Get Query (Fetch Single Item)

This example executes the `petrock_example_feature_name/get` query to retrieve a specific item by its ID. Remember to URL-encode the ID if it contains special characters (like spaces).

```bash
# Replace 'My%20First%20Item' with the URL-encoded ID of the item you want to fetch
curl "http://localhost:8080/queries/petrock_example_feature_name/get?ID=My%20First%20Item"
```

*(Expected Output: A JSON object representing the requested item, e.g., `{"id":"My First Item","name":"My First Item",...}`)*
