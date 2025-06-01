# MCP Server Support for Petrock Applications

## Overview

Add Model Context Protocol (MCP) server support to petrock-generated applications. When an application is called `blog`, the command `blog mcp` will start an MCP server on stdio using JSON-RPC 2.0.

## Goals

1. Add MCP server command to all petrock-generated applications
2. Implement a lightweight JSON-RPC 2.0 library in Go (no external dependencies)
3. Provide basic MCP server functionality in `core/mcp.go`
4. Enable applications to expose their data as MCP resources and tools

## Architecture

```
petrock-generated-app
├── cmd/
│   └── app_name/
│       ├── main.go      # registers mcp command
│       ├── serve.go     # existing HTTP server
│       └── mcp.go       # new MCP server command
└── core/
    ├── app.go           # existing core app
    ├── jsonrpc.go       # JSON-RPC 2.0 implementation
    └── mcp.go           # MCP server implementation
```

## Implementation Plan

### Phase 1: JSON-RPC 2.0 Foundation

1. **Create `core/jsonrpc.go`** - Simple JSON-RPC 2.0 implementation
   - Request/Response types
   - Notification handling
   - Error responses
   - Stdio transport layer

### Phase 2: MCP Core Implementation  

2. **Create `core/mcp.go`** - MCP server implementation
   - Server capability negotiation
   - Initialize/Ping handlers
   - Resource discovery and retrieval
   - Tool listing and execution
   - Basic logging support

### Phase 3: Command Integration

3. **Create `cmd/petrock_example_project_name/mcp.go`** - MCP command
   - Cobra command setup
   - Server initialization 
   - Stdio connection handling

4. **Update `cmd/petrock_example_project_name/main.go`** - Register MCP command

### Phase 4: Default Resources & Tools

5. **Implement default MCP features**
   - Database query tool
   - Key-value store access
   - Application info resources
   - Feature introspection

## File Structure

- `plan.md` - This plan document
- `jsonrpc-example.go` - JSON-RPC 2.0 implementation example
- `mcp-server-example.go` - MCP server core implementation example  
- `mcp-command-example.go` - Cobra command implementation example

## MCP Server Capabilities

The server will provide:

### Resources
- `app://info` - Application metadata and version
- `app://features` - List of registered features
- `app://config` - Application configuration (sanitized)

### Tools  
- `query_database` - Execute read-only SQL queries
- `get_kv` - Retrieve key-value store entries
- `list_kv_keys` - List available keys in KV store

### Prompts
- `analyze_app` - Template for analyzing the application structure
- `debug_feature` - Template for debugging specific features

## Testing Strategy

1. **Unit tests** for JSON-RPC implementation
2. **Integration tests** for MCP protocol compliance
3. **End-to-end tests** with actual MCP clients
4. **Manual testing** with Claude Desktop or other MCP clients

## Future Enhancements

- Custom resources/tools registration API
- Feature-specific MCP extensions  
- WebSocket transport support
- Advanced security/sandboxing
- Performance monitoring tools

## Dependencies

- Standard library only (no external dependencies)
- Existing petrock core functionality
- Cobra CLI framework (already used)
