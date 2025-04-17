# Technical Implementation Plan for Pointer-Only Commands

## Overview

The current implementation of Petrock's command pattern supports both value types and pointer types, leading to duplicated handler logic and potential inconsistencies. This plan outlines a migration to a pointer-only command pattern, which will simplify the codebase, reduce duplication, and establish a consistent approach to command handling.

## Detailed Task Breakdown

### T1: Core Command Interface Updates - DONE

**T1.1:** Update core.Command interface documentation - DONE
- Modify documentation to specify that commands should be pointer types - DONE
- Add examples showing proper command creation and usage - DONE

**T1.2:** Update core.Executor implementation - DONE
- Modify Execute method to expect pointer commands - DONE
- Add runtime checks to validate commands are pointers - DONE
- Add warning logs for non-pointer commands - DONE

**T1.3:** Update CommandRegistry implementation - DONE
- Modify Register method to accept only pointer type commands - DONE
- Update handler lookup logic to work exclusively with pointer types - DONE

**Definition of Done for T1:**
- Command interface documentation clearly states pointer requirement
- Executor rejects or warns about non-pointer commands
- CommandRegistry only accepts pointer types for registration

### T2: Feature Template Modifications - DONE

**T2.1:** Update commands.go template - DONE
- Modify type assertions to use pointer interface - DONE
- Update command registration examples to use pointers - DONE
- Update method receivers to use pointer types - DONE

**T2.2:** Update state.go template - DONE
- Modify Apply method to accept pointer commands only - DONE
- Update type switches to handle pointer variants only - DONE
- Update documentation examples - DONE

**T2.3:** Update execute.go template - DONE
- Simplify handler methods to only handle pointer types - DONE
- Remove duplicate handling code for value types - DONE
- Update error handling to match pointer-only approach - DONE

**T2.4:** Update worker.go template - DONE
- Modify command handling to use pointer commands only - DONE
- Ensure all command creations use pointer syntax - DONE

**Definition of Done for T2: - DONE**
- Feature template has been updated to use pointer-only commands - DONE
- No duplicate handling code for value types exists - DONE
- All command creations use consistent pointer syntax - DONE

### T3: Testing Infrastructure - SKIPPED

**T3.1:** Create migration test - SKIPPED
- Create test that verifies pointer commands work correctly - SKIPPED
- Add test cases for all command types - SKIPPED

**T3.2:** Update existing tests - SKIPPED
- Modify any tests that use value commands - SKIPPED
- Ensure all test assertions work with pointer types - SKIPPED

**Definition of Done for T3: - SKIPPED**
- Tests pass with pointer-only command usage - SKIPPED
- No test failures related to command type handling - SKIPPED

### T4: CLI Tool Updates

**T4.1:** Update petrock feature command
- Modify generated code to use pointer commands
- Update any examples in generated documentation

**T4.2:** Update code generation templates
- Ensure all command registration uses pointer syntax
- Update handler implementations to use pointer-only approach

**Definition of Done for T4:**
- CLI generates code using pointer-only commands
- Generated code documentation reflects pointer usage

### T5: Documentation Updates

**T5.1:** Update API documentation
- Modify command-related documentation to specify pointer usage
- Update examples in markdown files

**T5.2:** Create migration guide
- Document steps for users to migrate existing code
- Provide examples showing before/after changes

**Definition of Done for T5:**
- All documentation consistently shows pointer command usage
- Migration guide exists and clearly explains the transition

### T6: Migration Strategy

**T6.1:** Add backward compatibility layer
- Create temporary compatibility layer to handle value commands
- Add deprecation warnings when value commands are detected

**T6.2:** Implementation plan for existing projects
- Create upgrade script to help convert existing code
- Document breaking changes and necessary manual updates

**Definition of Done for T6:**
- Compatibility layer allows gradual adoption
- Clear path exists for users to migrate existing projects

## Implementation Details

### Command Registration

Current pattern:
```go
app.CommandRegistry.Register(CreateCommand{}, featureExecutor.HandleCreate, featureExecutor)
```

New pattern:
```go
app.CommandRegistry.Register(&CreateCommand{}, featureExecutor.HandleCreate, featureExecutor)
```

### Handler Simplification

Current pattern:
```go
switch cmd := command.(type) {
case CreateCommand:
    // Handle value type
    return e.state.Apply(cmd, msg)
case *CreateCommand:
    // Handle pointer type
    return e.state.Apply(*cmd, msg)
default:
    // Error handling
}
```

New pattern:
```go
switch cmd := command.(type) {
case *CreateCommand:
    // Handle pointer type only
    return e.state.Apply(cmd, msg)
default:
    // Error handling
}
```

### Command Creation

All command creation will use pointer syntax:

```go
cmd := &CreateCommand{
    Name: "example",
    Description: "Example command",
}

err := executor.Execute(ctx, cmd)
```

### State.Apply Method

Update the Apply method to accept pointer commands:

```go
func (s *State) Apply(payload interface{}, msg *core.Message) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    switch cmd := payload.(type) {
    case *CreateCommand:
        // Handle command
    case *UpdateCommand:
        // Handle command
    // Other commands
    }
    
    return nil
}
```

## Benefits

1. **Simplified code** - Removes duplicate handler logic for value and pointer types
2. **Reduced memory usage** - Avoids unnecessary copying of command structs
3. **Consistent pattern** - Establishes a single way to handle commands throughout the codebase
4. **Improved performance** - Reduces type assertions and unnecessary conversions
5. **Better IDE support** - Clearer type signatures improve code completion and analysis
6. **Reduced complexity** - Makes the codebase easier to maintain and understand