# CLI Tool

The Petrock CLI provides a command-line interface for working with Petrock applications. It handles project creation, feature generation, server management, and deployment.

## CLI Structure

```go
type CLI struct {
    app *cli.App
}

func NewCLI() *CLI {
    app := &cli.App{
        Name:  "petrock",
        Usage: "A command-line tool for Petrock applications",
        Commands: []*cli.Command{
            // Commands will be defined below
        },
    }
    
    return &CLI{app: app}
}

func (c *CLI) Run(args []string) error {
    return c.app.Run(args)
}
```

## Project Creation

```go
// Project creation command
var NewCommand = &cli.Command{
    Name:  "new",
    Usage: "Create a new Petrock application",
    Action: func(c *cli.Context) error {
        if c.NArg() != 1 {
            return fmt.Errorf("usage: petrock new [name]")
        }
        
        name := c.Args().First()
        return CreateNewProject(name)
    },
}

// CreateNewProject creates a new project structure
func CreateNewProject(name string) error {
    // Implementation...
    return nil
}
```

## Feature Generation

```go
// Feature generation command
var FeatureCommand = &cli.Command{
    Name:  "feature",
    Usage: "Generate a new feature module",
    Action: func(c *cli.Context) error {
        if c.NArg() != 1 {
            return fmt.Errorf("usage: petrock feature [name]")
        }
        
        name := c.Args().First()
        return CreateFeatureModule(".", name)
    },
}

// CreateFeatureModule generates a new feature module
func CreateFeatureModule(projectPath, featureName string) error {
    // Implementation...
    return nil
}
```

## Server Management

```go
// Server management command
var StartCommand = &cli.Command{
    Name:  "start",
    Usage: "Start the development server",
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:  "port",
            Value: "3000",
            Usage: "Port to listen on",
        },
    },
    Action: func(c *cli.Context) error {
        port := c.String("port")
        return StartServer(".", port)
    },
}

// StartServer starts the development server
func StartServer(projectPath, port string) error {
    // Implementation...
    return nil
}
```

## Build Command

```go
// Build command
var BuildCommand = &cli.Command{
    Name:  "build",
    Usage: "Build a production binary",
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:  "output",
            Value: "",
            Usage: "Output file name",
        },
        &cli.StringFlag{
            Name:  "os",
            Value: "",
            Usage: "Target OS (linux, darwin, windows)",
        },
        &cli.StringFlag{
            Name:  "arch",
            Value: "",
            Usage: "Target architecture (amd64, arm64)",
        },
    },
    Action: func(c *cli.Context) error {
        output := c.String("output")
        if output == "" {
            output = filepath.Base(filepath.Dir(".")) // Use directory name
        }
        
        return BuildBinary(".", output, c.String("os"), c.String("arch"))
    },
}

// BuildBinary builds a production binary
func BuildBinary(projectPath, outputPath, targetOS, targetArch string) error {
    // Implementation...
    return nil
}
```

## Deploy Command

```go
// Deploy command
var DeployCommand = &cli.Command{
    Name:  "deploy",
    Usage: "Deploy the application",
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:  "target",
            Value: "",
            Usage: "Deployment target (e.g., user@server:/path)",
            Required: true,
        },
    },
    Action: func(c *cli.Context) error {
        target := c.String("target")
        return Deploy(".", target)
    },
}

// Deploy deploys the application
func Deploy(projectPath, target string) error {
    // Implementation...
    return nil
}
```

## Database Commands

```go
// Database commands
var DBCommand = &cli.Command{
    Name:  "db",
    Usage: "Database management commands",
    Subcommands: []*cli.Command{
        {
            Name:  "backup",
            Usage: "Create a database backup",
            Action: func(c *cli.Context) error {
                return BackupDatabase(".")
            },
        },
        {
            Name:  "restore",
            Usage: "Restore a database backup",
            Action: func(c *cli.Context) error {
                if c.NArg() != 1 {
                    return fmt.Errorf("usage: petrock db restore [backup-file]")
                }
                
                backupFile := c.Args().First()
                return RestoreDatabase(".", backupFile)
            },
        },
        {
            Name:  "schema",
            Usage: "Print the database schema",
            Action: func(c *cli.Context) error {
                return PrintDatabaseSchema(".")
            },
        },
    },
}

// Database management functions
func BackupDatabase(projectPath string) error {
    // Implementation...
    return nil
}

func RestoreDatabase(projectPath, backupFile string) error {
    // Implementation...
    return nil
}

func PrintDatabaseSchema(projectPath string) error {
    // Implementation...
    return nil
}
```

## Generate Commands

```go
// Generate commands
var GenerateCommand = &cli.Command{
    Name:  "generate",
    Usage: "Generate various components",
    Subcommands: []*cli.Command{
        {
            Name:  "command",
            Usage: "Generate a new command",
            Action: func(c *cli.Context) error {
                if c.NArg() < 2 {
                    return fmt.Errorf("usage: petrock generate command [module] [name]")
                }
                
                module := c.Args().Get(0)
                name := c.Args().Get(1)
                return GenerateCommand(".", module, name)
            },
        },
        {
            Name:  "form",
            Usage: "Generate a new form",
            Action: func(c *cli.Context) error {
                if c.NArg() < 2 {
                    return fmt.Errorf("usage: petrock generate form [module] [name]")
                }
                
                module := c.Args().Get(0)
                name := c.Args().Get(1)
                return GenerateForm(".", module, name)
            },
        },
        {
            Name:  "job",
            Usage: "Generate a new background job",
            Action: func(c *cli.Context) error {
                if c.NArg() < 2 {
                    return fmt.Errorf("usage: petrock generate job [module] [name]")
                }
                
                module := c.Args().Get(0)
                name := c.Args().Get(1)
                return GenerateJob(".", module, name)
            },
        },
    },
}

// Generate component functions
func GenerateCommand(projectPath, module, name string) error {
    // Implementation...
    return nil
}

func GenerateForm(projectPath, module, name string) error {
    // Implementation...
    return nil
}

func GenerateJob(projectPath, module, name string) error {
    // Implementation...
    return nil
}
```

## List Command

```go
// List command
var ListCommand = &cli.Command{
    Name:  "list",
    Usage: "List various components",
    Subcommands: []*cli.Command{
        {
            Name:  "modules",
            Usage: "List all modules",
            Action: func(c *cli.Context) error {
                return ListModules(".")
            },
        },
        {
            Name:  "commands",
            Usage: "List all commands",
            Action: func(c *cli.Context) error {
                return ListCommands(".")
            },
        },
        {
            Name:  "routes",
            Usage: "List all routes",
            Action: func(c *cli.Context) error {
                return ListRoutes(".")
            },
        },
    },
}

// List component functions
func ListModules(projectPath string) error {
    // Implementation...
    return nil
}

func ListCommands(projectPath string) error {
    // Implementation...
    return nil
}

func ListRoutes(projectPath string) error {
    // Implementation...
    return nil
}
```

## SystemD Commands

```go
// SystemD commands
var SystemDCommand = &cli.Command{
    Name:  "systemd",
    Usage: "Generate and manage systemd service files",
    Subcommands: []*cli.Command{
        {
            Name:  "generate",
            Usage: "Generate a systemd service file",
            Action: func(c *cli.Context) error {
                return GenerateSystemDService(".")
            },
        },
        {
            Name:  "install",
            Usage: "Install a systemd service file",
            Action: func(c *cli.Context) error {
                return InstallSystemDService(".")
            },
        },
    },
}

// SystemD management functions
func GenerateSystemDService(projectPath string) error {
    // Implementation...
    return nil
}

func InstallSystemDService(projectPath string) error {
    // Implementation...
    return nil
}
```

## Project Structure Generation

The CLI tool generates the following project structure when creating a new project:

```
myapp/
├── main.go                 # Application entrypoint
├── app/
│   ├── shared/
│   │   └── ui.go           # App-specific shared components
│   └── .gitkeep
├── config/
│   └── app.go              # Application configuration
└── tmp/
    └── .gitkeep            # Temporary files
```

The main.go file would be generated with:

```go
// main.go content
const mainGoTemplate = `package main

import (
	"log"
	
	"github.com/yourusername/petrock/pkg/core"
)

func main() {
	app := core.NewApplication()
	
	if err := app.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
`
```

## Feature Module Generation

When generating a new feature module, the CLI creates the following files:

```
app/feature_name/
├── messages.go     # Command definitions
├── actions.go      # Command handlers
├── state.go        # State management
├── ui.go           # UI components
├── routes.go       # HTTP routes
└── main.go         # Module initialization
```

Each file is generated with appropriate templates, for example:

```go
// messages.go template
const messagesTemplate = `package %s

import (
	"time"
	
	"github.com/yourusername/petrock/pkg/core"
)

// Define your commands here
type Create%s struct {
	ID        string    ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	// Add your fields here
}

func (c Create%s) Type() string {
	return "%s.create"
}

func (c Create%s) EntityID() string {
	return c.ID
}

// Register commands with core
func RegisterCommands() {
	core.RegisterCommand(Create%s{})
	// Register other commands
}
`
```
