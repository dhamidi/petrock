# Migration Plan: Go Workspaces

## Overview

Migrate from custom build.sh orchestration to Go workspaces for simpler development workflow while maintaining template functionality.

## Current State

- `build.sh` copies `go.mod.skel` to `go.mod`, builds skeleton, then removes `go.mod` for embedding
- Skeleton uses placeholder imports `github.com/petrock/example_module_path` that don't resolve
- Two separate build contexts: petrock tool and skeleton template
- Complex build orchestration with potential for human error

## Target State

- Go workspace with both modules declared
- Skeleton uses replace directive to make placeholder imports resolve locally
- Standard `go build` commands work from any location
- Simplified build process with built-in Go tooling
- Same template functionality via string substitution during generation

## Migration Steps

### Phase 1: Create Workspace Structure

1. **Create go.work in repository root:**
   ```go
   go 1.23
   
   use (
       .
       ./internal/skeleton
   )
   ```

2. **Update internal/skeleton/go.mod.skel:**
   ```go
   module github.com/petrock/example_module_path
   
   go 1.23
   
   replace github.com/petrock/example_module_path => .
   
   require (
       github.com/mattn/go-sqlite3 v1.14.27
       github.com/spf13/cobra v1.8.1
       maragu.dev/gomponents v1.0.0
   )
   
   require (
       github.com/inconshreveable/mousetrap v1.1.0 // indirect
       github.com/spf13/pflag v1.0.5 // indirect
   )
   ```

3. **Copy go.mod.skel to go.mod in skeleton:**
   ```bash
   cp internal/skeleton/go.mod.skel internal/skeleton/go.mod
   ```

### Phase 2: Verify Build Compatibility

4. **Test skeleton builds independently:**
   ```bash
   cd internal/skeleton
   go build ./...
   go test ./...
   ```

5. **Test petrock builds from root:**
   ```bash
   go build ./cmd/petrock
   ```

6. **Test go:embed still works:**
   - Verify petrock binary includes skeleton files
   - Test project generation with `./petrock init test-project`
   - Verify generated code compiles after string substitution

### Phase 3: Simplify Build Process

7. **Update build.sh to remove complexity:**
   ```bash
   build_skeleton() {
     cd internal/skeleton
     go test ./...
     go build ./...
   }
   
   build_petrock() {
     go build ./cmd/...
     go install ./cmd/petrock
   }
   ```

8. **Remove go.mod copying/deletion logic**

9. **Update AGENT.md commands:**
   - Build: `go build ./cmd/petrock`
   - Test skeleton: `cd internal/skeleton && go build ./...`
   - Full build: `./build.sh` (simplified)

### Phase 4: Integration Testing

10. **Run full integration test suite:**
    ```bash
    ./test-all.sh
    ```

11. **Test generation pipeline:**
    - Generate new project: `./petrock init test-workspace`
    - Add feature: `./petrock feature posts`
    - Verify builds: `cd test-workspace && go build ./...`
    - Start server: `cd test-workspace && go run ./cmd/test-workspace serve`

12. **Test UI gallery access:**
    - Navigate to `http://localhost:8080/_/ui`
    - Verify component pages load correctly

## Testing Checklist

- [ ] `go build ./cmd/petrock` works from root
- [ ] `cd internal/skeleton && go build ./...` works
- [ ] `cd internal/skeleton && go test ./...` passes
- [ ] `./petrock init test-project` generates working project
- [ ] Generated project compiles: `cd test-project && go build ./...`
- [ ] `./petrock feature posts` adds working feature
- [ ] Feature integration compiles and runs
- [ ] UI gallery accessible and functional
- [ ] All integration tests pass: `./test-all.sh`

## Rollback Plan

If migration fails:

1. **Remove workspace files:**
   ```bash
   rm go.work
   rm internal/skeleton/go.mod
   ```

2. **Restore original build.sh logic**

3. **Verify original build process works:**
   ```bash
   ./build.sh
   ```

## Benefits

- **Simplified development**: Standard Go commands work everywhere
- **IDE compatibility**: Better code navigation and completion
- **Reduced build complexity**: No custom orchestration needed
- **Maintainability**: Uses Go's built-in workspace features
- **Developer experience**: Familiar workflow for Go developers

## Risks & Mitigations

**Risk**: String substitution breaks with replace directive
- **Mitigation**: Replace directive is local to skeleton module, generation still does string replacement on file contents

**Risk**: go:embed behavior changes
- **Mitigation**: Embed paths remain the same, workspace doesn't affect embedding

**Risk**: Template validation fails
- **Mitigation**: Same compilation occurs, just with workspace resolution instead of isolation

**Risk**: Generated projects reference workspace
- **Mitigation**: String substitution replaces all placeholder imports during generation

## Success Criteria

- Both build commands work without build.sh orchestration
- Full template generation pipeline functions identically  
- Integration tests pass
- UI gallery remains functional
- Developer workflow simplified
- Build times maintained or improved

## Timeline

- **Phase 1-2**: 30 minutes (setup and verification)
- **Phase 3**: 15 minutes (build simplification)  
- **Phase 4**: 30 minutes (comprehensive testing)
- **Total**: ~75 minutes

## Post-Migration

- Update documentation in README.md
- Remove obsolete build.sh complexity
- Consider removing .skel extension from go.mod.skel
- Update CI/CD if applicable
