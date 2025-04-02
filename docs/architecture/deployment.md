# Deployment

Petrock applications are designed for simple, straightforward deployment. The Petrock framework provides tools and utilities for building, deploying, and managing applications in production environments.

## Deployment Strategy

Petrock uses a simple deployment strategy:

1. Build a static binary with all dependencies included
2. Transfer the binary to the target server
3. Configure the application on the target server
4. Start the application (optionally using systemd or another supervisor)

## Binary Building

Applications are compiled into a single static binary:

```go
// BuildBinary builds a production binary
func BuildBinary(projectPath, outputPath, targetOS, targetArch string) error {
    // Set up build environment
    env := []string{
        "CGO_ENABLED=1", // Required for SQLite
        "GO111MODULE=on",
    }
    
    // Set target OS/Arch if specified
    if targetOS != "" {
        env = append(env, "GOOS="+targetOS)
    }
    
    if targetArch != "" {
        env = append(env, "GOARCH="+targetArch)
    }
    
    // Build command
    cmd := exec.Command("go", "build", "-o", outputPath)
    cmd.Dir = projectPath
    cmd.Env = append(os.Environ(), env...)
    
    // Capture output
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("build failed: %v\n%s", err, output)
    }
    
    return nil
}
```

## Cross-Compilation

Petrock supports cross-compilation for different target platforms:

```go
// Cross-compile for different platforms
func CrossCompile(projectPath string) error {
    targets := []struct {
        os   string
        arch string
    }{
        {"linux", "amd64"},
        {"linux", "arm64"},
        {"darwin", "amd64"},
        {"darwin", "arm64"},
        {"windows", "amd64"},
    }
    
    for _, target := range targets {
        outputName := fmt.Sprintf("petrock-app-%s-%s", target.os, target.arch)
        if target.os == "windows" {
            outputName += ".exe"
        }
        
        if err := BuildBinary(projectPath, outputName, target.os, target.arch); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Deployment

Deployment transfers the binary to the target server:

```go
// Deploy the application to a server
func Deploy(projectPath, target string) error {
    // Build the binary
    binaryPath := filepath.Join(projectPath, "petrock-app")
    if err := BuildBinary(projectPath, binaryPath, "", ""); err != nil {
        return err
    }
    
    // Transfer the binary to the target server
    cmd := exec.Command("rsync", "-avz", binaryPath, target)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("deploy failed: %v\n%s", err, output)
    }
    
    return nil
}
```

## SystemD Service

Petrock can generate and install systemd service files:

```go
// Generate a systemd service file
func GenerateSystemDService(projectPath string) error {
    projectName := filepath.Base(projectPath)
    
    serviceContent := fmt.Sprintf(`[Unit]
Description=%s service
After=network.target

[Service]
Type=simple
User=%s
WorkingDirectory=/opt/%s
ExecStart=/opt/%s/petrock-app
Restart=on-failure

[Install]
WantedBy=multi-user.target
`, projectName, os.Getenv("USER"), projectName, projectName)
    
    return ioutil.WriteFile(filepath.Join(projectPath, projectName+".service"), []byte(serviceContent), 0644)
}

// Install a systemd service file
func InstallSystemDService(projectPath string) error {
    projectName := filepath.Base(projectPath)
    servicePath := filepath.Join(projectPath, projectName+".service")
    
    // Generate the service file if it doesn't exist
    if _, err := os.Stat(servicePath); os.IsNotExist(err) {
        if err := GenerateSystemDService(projectPath); err != nil {
            return err
        }
    }
    
    // Copy to systemd directory
    cmd := exec.Command("sudo", "cp", servicePath, "/etc/systemd/system/")
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Reload systemd
    cmd = exec.Command("sudo", "systemctl", "daemon-reload")
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Enable the service
    cmd = exec.Command("sudo", "systemctl", "enable", projectName)
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Start the service
    cmd = exec.Command("sudo", "systemctl", "start", projectName)
    if err := cmd.Run(); err != nil {
        return err
    }
    
    return nil
}
```

## Configuration

Petrock applications use environment variables for configuration:

```go
// Configuration structure
type Config struct {
    Port        string `env:"PORT" default:"3000"`
    Environment string `env:"ENVIRONMENT" default:"development"`
    DatabaseURL string `env:"DATABASE_URL" default:"petrock.db"`
    LogLevel    string `env:"LOG_LEVEL" default:"info"`
    SecretKey   string `env:"SECRET_KEY" required:"true"`
}

// Load configuration from environment
func LoadConfig() (*Config, error) {
    var config Config
    
    if err := env.Parse(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## Database Management

Petrock includes utilities for database management:

```go
// Backup the database
func BackupDatabase(dbPath, backupPath string) error {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }
    defer db.Close()
    
    backup, err := sql.Open("sqlite3", backupPath)
    if err != nil {
        return err
    }
    defer backup.Close()
    
    // Perform backup
    _, err = db.Exec("VACUUM INTO ?", backupPath)
    return err
}

// Restore a database from backup
func RestoreDatabase(backupPath, dbPath string) error {
    // Simple file copy for SQLite
    input, err := ioutil.ReadFile(backupPath)
    if err != nil {
        return err
    }
    
    return ioutil.WriteFile(dbPath, input, 0644)
}
```

## Monitoring

Petrock provides built-in monitoring endpoints:

```go
// Register monitoring endpoints
func RegisterMonitoringEndpoints(router *Router) {
    router.Get("/healthz", HandleHealthCheck)
    router.Get("/readyz", HandleReadinessCheck)
    router.Get("/metrics", HandleMetrics)
}

// Health check handler
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}

// Readiness check handler
func HandleReadinessCheck(w http.ResponseWriter, r *http.Request) {
    // Check database connection
    db := core.GetDB()
    err := db.Ping()
    
    w.Header().Set("Content-Type", "application/json")
    
    if err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte(`{"status":"not_ready","reason":"database_unavailable"}`))
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ready"}`))
}

// Metrics handler
func HandleMetrics(w http.ResponseWriter, r *http.Request) {
    // Collect basic metrics
    metrics := map[string]interface{}{
        "uptime":        time.Since(startTime).Seconds(),
        "goroutines":    runtime.NumGoroutine(),
        "memory":        getMemoryStats(),
        "requests":      getRequestStats(),
        "message_count": getMessageCount(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(metrics)
}
```

## Deployment Checklist

Petrock includes a deployment checklist command:

```go
// Deployment checklist command
var ChecklistCommand = &cli.Command{
    Name:  "checklist",
    Usage: "Run pre-deployment checks",
    Action: func(c *cli.Context) error {
        return RunDeploymentChecklist(".")
    },
}

// Run deployment checklist
func RunDeploymentChecklist(projectPath string) error {
    checks := []struct {
        name     string
        check    func(string) error
        critical bool
    }{
        {"Configuration", CheckConfiguration, true},
        {"Database", CheckDatabase, true},
        {"Build", CheckBuild, true},
        {"Routes", CheckRoutes, false},
        {"Static Files", CheckStaticFiles, false},
    }
    
    fmt.Println("Running deployment checklist:")
    
    var failed bool
    
    for _, check := range checks {
        fmt.Printf("  - %s: ", check.name)
        
        if err := check.check(projectPath); err != nil {
            fmt.Printf("FAILED (%s)\n", err)
            if check.critical {
                failed = true
            }
        } else {
            fmt.Println("PASSED")
        }
    }
    
    if failed {
        return errors.New("critical checks failed")
    }
    
    return nil
}
```

## Environment Setup

Functions for setting up the deployment environment:

```go
// Setup deployment environment
func SetupEnvironment(target string) error {
    cmds := []struct {
        name string
        args []string
    }{
        {"Create application directory", []string{"mkdir", "-p", "/opt/petrock-app"}},
        {"Set permissions", []string{"chown", "-R", "petrock:petrock", "/opt/petrock-app"}},
        {"Create log directory", []string{"mkdir", "-p", "/var/log/petrock-app"}},
        {"Set log permissions", []string{"chown", "-R", "petrock:petrock", "/var/log/petrock-app"}},
    }
    
    for _, c := range cmds {
        fmt.Printf("Running: %s\n", c.name)
        cmd := exec.Command("ssh", target, strings.Join(c.args, " "))
        if err := cmd.Run(); err != nil {
            return fmt.Errorf("%s failed: %v", c.name, err)
        }
    }
    
    return nil
}
```

## Zero-Downtime Deployment

Petrock supports zero-downtime deployments using a simple strategy:

```go
// Perform a zero-downtime deployment
func ZeroDowntimeDeployment(projectPath, target string) error {
    // Build the binary
    binaryPath := filepath.Join(projectPath, "petrock-app")
    if err := BuildBinary(projectPath, binaryPath, "", ""); err != nil {
        return err
    }
    
    // Transfer the binary to a temporary location on the target server
    tempPath := "/opt/petrock-app/petrock-app.new"
    cmd := exec.Command("rsync", "-avz", binaryPath, target+":"+tempPath)
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Make the new binary executable
    cmd = exec.Command("ssh", target, "chmod", "+x", tempPath)
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Replace the old binary with the new one
    cmd = exec.Command("ssh", target, "mv", tempPath, "/opt/petrock-app/petrock-app")
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Send SIGHUP to the running process
    cmd = exec.Command("ssh", target, 
        "systemctl", "reload", "petrock-app")
    if err := cmd.Run(); err != nil {
        return err
    }
    
    return nil
}
```
