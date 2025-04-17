# Technical Implementation Plan for Worker Self-Inspection

## Overview

Currently, the self-inspect command returns information about commands, queries, routes, and features, but does not include details about registered workers. This plan outlines the steps needed to enhance the inspection capabilities to include worker information, making it easier for developers to understand the background processing capabilities of their application.

## Detailed Task Breakdown

### T1: Core Inspection Interface Updates

**T1.1:** Update InspectResult structure
- [x] Add a Workers field to InspectResult struct in core/inspect.go
- [x] Define a WorkerSchema struct to represent worker metadata

**T1.2:** Create schema builder for workers
- [x] Implement buildWorkerSchema function to extract worker metadata
- [x] Define fields to capture in worker schema (name, type, methods)

**T1.3:** Update GetInspectResult implementation
- [x] Modify to gather worker information from App.workers
- [x] Convert workers to schema representation

**Definition of Done for T1:** ✅
- ✅ InspectResult includes worker information
- ✅ Worker metadata is properly extracted using reflection
- ✅ GetInspectResult returns complete worker information

### T2: Worker Interface Enhancements

**T2.1:** Add name/description capability to Worker interface
- [x] Add optional WorkerInfo() method to Worker interface
- [x] Implement default fallback for workers without explicit info

**T2.2:** Update existing worker implementations
- [x] Add WorkerInfo() to feature template worker
- [x] Ensure worker type information is accessible via reflection

**Definition of Done for T2:** ✅
- ✅ Workers can provide self-descriptive information
- ✅ Fallback mechanisms exist for workers without explicit descriptions

### T3: Test Infrastructure Updates

**T3.1:** Update test command
- [x] Modify cmd/petrock/test.go to verify worker information
- [x] Add assertions for worker fields in self-inspect output

**T3.2:** Create worker inspection tests
- [x] Add test cases for worker self-description
- [x] Verify reflection-based extraction works correctly

**Definition of Done for T3:** ✅
- ✅ Tests verify worker information is correctly included
- ✅ Test coverage for worker inspection functionality

### T4: Documentation Updates

**T4.1:** Update API documentation
- [x] Document worker inspection capabilities
- [x] Add examples showing how to query worker information

**T4.2:** Update self-inspect command documentation
- [x] Describe worker fields in output
- [x] Provide examples of worker information usage

**Definition of Done for T4:** ✅
- ✅ Documentation accurately reflects worker inspection capabilities
- ✅ Examples demonstrate how to use the feature

## Implementation Details

### Worker Schema Structure

```go
// WorkerSchema represents metadata about a registered worker
type WorkerSchema struct {
    Name        string   `json:"name"`        // Worker name or type
    Description string   `json:"description"` // Worker description if available
    Type        string   `json:"type"`        // Go type name
    Methods     []string `json:"methods"`     // Available methods
}
```

### InspectResult Updates

```go
// InspectResult holds application metadata
type InspectResult struct {
    Commands []CommandSchema `json:"commands"` // Schema of all registered commands
    Queries  []QuerySchema  `json:"queries"`  // Schema of all registered queries
    Routes   []string       `json:"routes"`   // List of all registered HTTP routes
    Features []string       `json:"features"` // List of all registered features
    Workers  []WorkerSchema `json:"workers"`  // Schema of all registered workers
}
```

### Worker Interface Enhancement

```go
// WorkerInfo provides optional self-description for workers
type WorkerInfo struct {
    Name        string
    Description string
}

// Worker is the interface that must be implemented by background workers
type Worker interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Work() error
    
    // Optional method for self-description
    // If not implemented, information will be extracted via reflection
    WorkerInfo() *WorkerInfo
}
```

### Worker Schema Extraction

```go
// buildWorkerSchema creates a schema from a worker instance
func buildWorkerSchema(worker Worker) WorkerSchema {
    schema := WorkerSchema{
        Type: fmt.Sprintf("%T", worker),
    }
    
    // Try to get WorkerInfo if implemented
    if infoProvider, ok := worker.(interface{ WorkerInfo() *WorkerInfo }); ok {
        if info := infoProvider.WorkerInfo(); info != nil {
            schema.Name = info.Name
            schema.Description = info.Description
        }
    }
    
    // Default name from type if not provided
    if schema.Name == "" {
        schema.Name = extractTypeName(schema.Type)
    }
    
    // Extract methods using reflection
    schema.Methods = extractWorkerMethods(worker)
    
    return schema
}
```

## Benefits

1. **Complete application introspection** - Adds workers to the self-inspection capabilities
2. **Improved developer experience** - Makes it easier to understand all components of an application
3. **Better documentation** - Self-documenting worker instances
4. **Debugging assistance** - Helps developers identify registered workers
5. **Consistency** - Completes the pattern of inspecting all major application components

## Implementation Complete

All planned tasks have been completed. The codebase now supports worker self-inspection throughout the system.

- u2705 Core worker interface updated with optional WorkerInfo method
- u2705 InspectResult structure updated to include worker information
- u2705 Worker schema generation implemented
- u2705 Feature template worker updated with WorkerInfo implementation
- u2705 Tests updated to verify worker information is included
- u2705 Documentation updated to reflect worker inspection capabilities