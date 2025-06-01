# Component Generator Development Plan

**FEATURE OBJECTIVE**: Implement individual component generators (`petrock new command|query|worker <feature-name>/<entity-name> [args]`) that extract and generate specific components from existing skeleton templates, with collision detection via `self inspect` subcommand integration.

**ACCEPTANCE CRITERIA**: 
- `petrock new command posts/create` generates only command files from skeleton template
- `petrock new query posts/get` generates only query files from skeleton template  
- `petrock new worker posts/summary` generates only worker files from skeleton template
- Generators detect existing components via `self inspect` and prevent overwriting
- Generated components integrate seamlessly with existing features
- All generators use existing string substitution mechanism with proper placeholder replacement
- Build commands `./build.sh` pass after component generation

## IMPLEMENTATION PHASES

### Phase 1: Core Generator Infrastructure
- **Task 1.1**: [File: cmd/petrock/new.go] Extend existing new command with component subcommands (Effort: Medium, Dependencies: None)
  - Types: ComponentType enum (command|query|worker), GenerateComponentOptions struct
  - Functions: NewComponentCmd(), validateComponentArgs(), parseFeatureEntityName()
  - Packages: github.com/spf13/cobra

- **Task 1.2**: [File: internal/generator/inspector.go] Create self-inspect integration for collision detection (Effort: Medium, Dependencies: 1.1)
  - Types: ComponentInspector interface, InspectResult struct
  - Functions: NewComponentInspector(), InspectExistingComponents(), ComponentExists()
  - Packages: os/exec, encoding/json

- **Task 1.3**: [File: internal/generator/component.go] Create component extraction and generation engine (Effort: High, Dependencies: 1.1, 1.2)
  - Types: ComponentGenerator interface, ComponentTemplate struct, ExtractionOptions struct
  - Functions: NewComponentGenerator(), ExtractComponent(), GenerateComponent()
  - Packages: path/filepath, internal/utils

### Phase 2: Command Generator Implementation
- **Task 2.1**: [File: internal/generator/command.go] Implement command-specific template extraction (Effort: Medium, Dependencies: 1.3)
  - Types: CommandGenerator struct implementing ComponentGenerator
  - Functions: ExtractCommandFiles(), GenerateCommandComponent(), validateCommandStructure()
  - Packages: internal/skeleton (embedded)

- **Task 2.2**: [File: internal/generator/templates/command.go] Define command file mapping and placeholders (Effort: Low, Dependencies: 2.1)
  - Types: CommandFileMap map[string]string, CommandPlaceholders struct  
  - Functions: GetCommandTemplateFiles(), GetCommandReplacements()
  - Packages: None

- **Task 2.3**: [File: cmd/petrock/new_command.go] Create command-specific CLI handler (Effort: Low, Dependencies: 2.1, 2.2)
  - Types: NewCommandOptions struct
  - Functions: NewCommandSubcommand(), runCommandGeneration(), validateCommandOptions()
  - Packages: github.com/spf13/cobra

### Phase 3: Query Generator Implementation  
- **Task 3.1**: [File: internal/generator/query.go] Implement query-specific template extraction (Effort: Medium, Dependencies: 1.3)
  - Types: QueryGenerator struct implementing ComponentGenerator
  - Functions: ExtractQueryFiles(), GenerateQueryComponent(), validateQueryStructure()
  - Packages: internal/skeleton (embedded)

- **Task 3.2**: [File: internal/generator/templates/query.go] Define query file mapping and placeholders (Effort: Low, Dependencies: 3.1) 
  - Types: QueryFileMap map[string]string, QueryPlaceholders struct
  - Functions: GetQueryTemplateFiles(), GetQueryReplacements()
  - Packages: None

- **Task 3.3**: [File: cmd/petrock/new_query.go] Create query-specific CLI handler (Effort: Low, Dependencies: 3.1, 3.2)
  - Types: NewQueryOptions struct  
  - Functions: NewQuerySubcommand(), runQueryGeneration(), validateQueryOptions()
  - Packages: github.com/spf13/cobra

### Phase 4: Worker Generator Implementation
- **Task 4.1**: [File: internal/generator/worker.go] Implement worker-specific template extraction (Effort: Medium, Dependencies: 1.3)
  - Types: WorkerGenerator struct implementing ComponentGenerator  
  - Functions: ExtractWorkerFiles(), GenerateWorkerComponent(), validateWorkerStructure()
  - Packages: internal/skeleton (embedded)

- **Task 4.2**: [File: internal/generator/templates/worker.go] Define worker file mapping and placeholders (Effort: Low, Dependencies: 4.1)
  - Types: WorkerFileMap map[string]string, WorkerPlaceholders struct
  - Functions: GetWorkerTemplateFiles(), GetWorkerReplacements()  
  - Packages: None

- **Task 4.3**: [File: cmd/petrock/new_worker.go] Create worker-specific CLI handler (Effort: Low, Dependencies: 4.1, 4.2)
  - Types: NewWorkerOptions struct
  - Functions: NewWorkerSubcommand(), runWorkerGeneration(), validateWorkerOptions()
  - Packages: github.com/spf13/cobra

### Phase 5: Template File Mapping System
- **Task 5.1**: [File: internal/generator/mappings.go] Create template-to-target file mapping system (Effort: Medium, Dependencies: 2.2, 3.2, 4.2)
  - Types: FileMapping struct, ComponentType enum, MappingConfig map[ComponentType]FileMapping
  - Functions: GetComponentFiles(), MapTemplateToTarget(), FilterComponentFiles()
  - Packages: path/filepath

- **Task 5.2**: [File: internal/generator/extractor.go] Implement selective file extraction from skeleton (Effort: High, Dependencies: 5.1)
  - Types: SkeletonExtractor struct, ExtractionFilter func(string) bool
  - Functions: NewSkeletonExtractor(), ExtractFiles(), ApplyFilter(), CopyFilteredFiles()
  - Packages: embed, io/fs, internal/utils

- **Task 5.3**: [File: internal/generator/placeholders.go] Centralize placeholder replacement logic (Effort: Medium, Dependencies: 5.2)
  - Types: PlaceholderReplacer struct, ReplacementRule struct
  - Functions: NewPlaceholderReplacer(), AddReplacement(), ProcessFile(), GetComponentReplacements()
  - Packages: strings, regexp

### Phase 6: Integration and Testing
- **Task 6.1**: [File: cmd/petrock/new.go] Integrate component subcommands into main new command (Effort: Low, Dependencies: 2.3, 3.3, 4.3)
  - Types: No new types
  - Functions: registerComponentSubcommands(), addComponentFlags()
  - Packages: github.com/spf13/cobra

- **Task 6.2**: [File: internal/generator/generator_test.go] Create comprehensive test suite for generators (Effort: High, Dependencies: All previous)
  - Types: GeneratorTestSuite struct, TestCase struct
  - Functions: TestCommandGeneration(), TestQueryGeneration(), TestWorkerGeneration(), TestCollisionDetection()
  - Packages: testing, os, path/filepath

- **Task 6.3**: [File: docs/plans/generators/examples.md] Create usage examples and documentation (Effort: Low, Dependencies: 6.1)
  - Types: No types (documentation)
  - Functions: No functions (documentation)
  - Packages: None

**CRITICAL PATH**: 
1. Task 1.1 (CLI infrastructure) → Task 1.3 (generation engine) → Task 5.2 (file extraction) → Component-specific generators (2.1, 3.1, 4.1) → Integration (6.1)

**TECHNICAL RISKS**:
- **Risk 1**: Template file interdependencies may require partial feature generation → **Mitigation**: Map component dependencies and generate required supporting files
- **Risk 2**: Self-inspect command execution may fail in incomplete projects → **Mitigation**: Gracefully handle inspect failures and provide fallback detection
- **Risk 3**: Placeholder replacement conflicts between components → **Mitigation**: Use component-specific placeholder scoping and validation
- **Risk 4**: Generated components may not compile without full feature context → **Mitigation**: Generate minimal supporting stubs and clear error messages

**BUILD REQUIREMENTS**:
- **Go Dependencies**: No new external dependencies required (uses existing cobra, embed)
- **Build Commands**: `./build.sh` to test full compilation pipeline, `go test ./internal/generator/...` for unit tests
- **Testing Strategy**: Integration tests generating sample components in temporary projects, validation via compilation and `self inspect`

**FIRST IMPLEMENTATION STEP**: 
Create the core CLI structure in `cmd/petrock/new.go` by adding component subcommands (`command`, `query`, `worker`) to the existing `new` command, using cobra's subcommand pattern established in the codebase.

## Component File Mappings

### Command Component Files
From `internal/skeleton/petrock_example_feature_name/`:
- `commands/base.go` → `{feature}/commands/base.go`
- `commands/{entity}.go` → `{feature}/commands/{entity}.go` 
- `commands/register.go` → `{feature}/commands/register.go`

### Query Component Files  
From `internal/skeleton/petrock_example_feature_name/`:
- `queries/base.go` → `{feature}/queries/base.go`
- `queries/{entity}.go` → `{feature}/queries/{entity}.go`

### Worker Component Files
From `internal/skeleton/petrock_example_feature_name/`:
- `workers/main.go` → `{feature}/workers/main.go`
- `workers/{entity}_worker.go` → `{feature}/workers/{entity}_worker.go`
- `workers/types.go` → `{feature}/workers/types.go`

## Placeholder Replacements

All generators will use existing petrock placeholders:
- `petrock_example_feature_name` → actual feature name (e.g., `posts`)
- `github.com/petrock/example_module_path` → target module path
- Template entity names → actual entity name (e.g., `create`, `get`, `summary`)

## Integration Points

1. **Self Inspect Integration**: Use `go run . self inspect --json` to detect existing components
2. **File System**: Generate into existing project structure, creating directories as needed  
3. **Registration**: Generated components must register themselves in existing registry systems
4. **Validation**: Ensure generated code compiles and integrates with existing features
