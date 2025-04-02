# Background Jobs

The background job system in Petrock allows for the execution of long-running tasks outside the HTTP request cycle. Jobs are persisted in the database and processed by worker goroutines.

## Job System Design

Jobs in Petrock follow a specific pattern:

1. Job builds state from the message log
2. Job runs in a separate goroutine
3. Before performing any action, the job catches up with the event log
4. Job performs actions based on its state
5. Job persists the fact that it performed an action in the message log

## Job Definitions

```go
// Job interface
type JobHandler func(ctx context.Context, params interface{}) error

// Job representation
type Job struct {
    ID        string
    Name      string
    Params    interface{}
    RunAt     time.Time
    Status    string // "pending", "running", "completed", "failed"
    CreatedAt time.Time
    UpdatedAt time.Time
    Error     string
}

// Job store interface
type JobStore interface {
    Save(job *Job) error
    Get(id string) (*Job, error)
    List(status string, limit int) ([]*Job, error)
    Update(job *Job) error
}
```

## Worker Implementation

```go
// Worker manages job execution
type Worker struct {
    store    JobStore
    handlers map[string]JobHandler
    stop     chan struct{}
    wg       sync.WaitGroup
}

func NewWorker(store JobStore) *Worker {
    return &Worker{
        store:    store,
        handlers: make(map[string]JobHandler),
        stop:     make(chan struct{}),
    }
}

func (w *Worker) RegisterHandler(name string, handler JobHandler) {
    w.handlers[name] = handler
}

func (w *Worker) Start(concurrency int) {
    for i := 0; i < concurrency; i++ {
        w.wg.Add(1)
        go w.processJobs()
    }
}

func (w *Worker) Stop() {
    close(w.stop)
    w.wg.Wait()
}

func (w *Worker) processJobs() {
    defer w.wg.Done()
    
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-w.stop:
            return
        case <-ticker.C:
            w.processNextBatch()
        }
    }
}

func (w *Worker) processNextBatch() {
    jobs, err := w.store.List("pending", 10)
    if err != nil {
        // Log error
        return
    }
    
    for _, job := range jobs {
        if job.RunAt.After(time.Now()) {
            continue
        }
        
        handler, ok := w.handlers[job.Name]
        if !ok {
            // Log unknown job type
            job.Status = "failed"
            job.Error = "Unknown job type"
            w.store.Update(job)
            continue
        }
        
        // Update job status
        job.Status = "running"
        job.UpdatedAt = time.Now()
        w.store.Update(job)
        
        // Run the job in a goroutine
        go func(j *Job) {
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
            defer cancel()
            
            err := handler(ctx, j.Params)
            
            j.UpdatedAt = time.Now()
            if err != nil {
                j.Status = "failed"
                j.Error = err.Error()
            } else {
                j.Status = "completed"
            }
            
            w.store.Update(j)
        }(job)
    }
}
```

## SQLite Job Store

```go
// SQLite implementation of JobStore
type SQLiteJobStore struct {
    db *sql.DB
    mu sync.RWMutex
}

func NewSQLiteJobStore(db *sql.DB) (*SQLiteJobStore, error) {
    // Create jobs table if it doesn't exist
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS jobs (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            params BLOB NOT NULL,
            run_at TIMESTAMP NOT NULL,
            status TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL,
            error TEXT
        );
        CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
        CREATE INDEX IF NOT EXISTS idx_jobs_run_at ON jobs(run_at);
    `)
    
    return &SQLiteJobStore{db: db}, err
}

func (s *SQLiteJobStore) Save(job *Job) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    params, err := json.Marshal(job.Params)
    if err != nil {
        return err
    }
    
    _, err = s.db.Exec(
        "INSERT INTO jobs (id, name, params, run_at, status, created_at, updated_at, error) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        job.ID,
        job.Name,
        params,
        job.RunAt,
        job.Status,
        job.CreatedAt,
        job.UpdatedAt,
        job.Error,
    )
    
    return err
}

func (s *SQLiteJobStore) Get(id string) (*Job, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var job Job
    var params []byte
    
    err := s.db.QueryRow(
        "SELECT id, name, params, run_at, status, created_at, updated_at, error FROM jobs WHERE id = ?",
        id,
    ).Scan(&job.ID, &job.Name, &params, &job.RunAt, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Error)
    
    if err != nil {
        return nil, err
    }
    
    // Deserialize params based on job type
    job.Params = deserializeJobParams(job.Name, params)
    
    return &job, nil
}

func (s *SQLiteJobStore) List(status string, limit int) ([]*Job, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var query string
    var args []interface{}
    
    if status == "" {
        query = "SELECT id, name, params, run_at, status, created_at, updated_at, error FROM jobs ORDER BY run_at ASC LIMIT ?"
        args = []interface{}{limit}
    } else {
        query = "SELECT id, name, params, run_at, status, created_at, updated_at, error FROM jobs WHERE status = ? ORDER BY run_at ASC LIMIT ?"
        args = []interface{}{status, limit}
    }
    
    rows, err := s.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var jobs []*Job
    
    for rows.Next() {
        var job Job
        var params []byte
        
        if err := rows.Scan(&job.ID, &job.Name, &params, &job.RunAt, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.Error); err != nil {
            return nil, err
        }
        
        // Deserialize params based on job type
        job.Params = deserializeJobParams(job.Name, params)
        
        jobs = append(jobs, &job)
    }
    
    return jobs, nil
}

func (s *SQLiteJobStore) Update(job *Job) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    params, err := json.Marshal(job.Params)
    if err != nil {
        return err
    }
    
    _, err = s.db.Exec(
        "UPDATE jobs SET name = ?, params = ?, run_at = ?, status = ?, updated_at = ?, error = ? WHERE id = ?",
        job.Name,
        params,
        job.RunAt,
        job.Status,
        job.UpdatedAt,
        job.Error,
        job.ID,
    )
    
    return err
}
```

## Job Registration

Jobs are registered with the job system:

```go
// Register job handlers
func RegisterJobs(worker *Worker) {
    worker.RegisterHandler("email.send", SendEmailJob)
    worker.RegisterHandler("image.process", ProcessImageJob)
    worker.RegisterHandler("report.generate", GenerateReportJob)
}

// Example job implementation
func SendEmailJob(ctx context.Context, params interface{}) error {
    emailParams := params.(EmailParams)
    
    // Build state from message log
    store := core.GetLogStore()
    lastVersion, _ := store.Version()
    
    // Send email
    err := sendEmail(emailParams.To, emailParams.Subject, emailParams.Body)
    if err != nil {
        return err
    }
    
    // Record that email was sent
    cmd := EmailSentCommand{
        EmailID:   emailParams.ID,
        To:        emailParams.To,
        Subject:   emailParams.Subject,
        SentAt:    time.Now(),
    }
    
    _, err = store.Append([]core.Message{cmd})
    return err
}
```

## Job Enqueuing

Jobs can be enqueued from anywhere in the application:

```go
// Create and enqueue a job
func EnqueueJob(store JobStore, name string, params interface{}, runAt time.Time) (string, error) {
    if runAt.IsZero() {
        runAt = time.Now()
    }
    
    job := &Job{
        ID:        core.NewID(),
        Name:      name,
        Params:    params,
        RunAt:     runAt,
        Status:    "pending",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    return job.ID, store.Save(job)
}

// Example usage
func ScheduleEmailSending(to, subject, body string) (string, error) {
    params := EmailParams{
        ID:      core.NewID(),
        To:      to,
        Subject: subject,
        Body:    body,
    }
    
    return EnqueueJob(jobStore, "email.send", params, time.Now())
}
```

## Job Status Checking

The status of jobs can be checked:

```go
// Check job status
func GetJobStatus(store JobStore, id string) (string, error) {
    job, err := store.Get(id)
    if err != nil {
        return "", err
    }
    
    return job.Status, nil
}

// Cancel a job
func CancelJob(store JobStore, id string) error {
    job, err := store.Get(id)
    if err != nil {
        return err
    }
    
    // Only pending jobs can be canceled
    if job.Status != "pending" {
        return fmt.Errorf("job is not pending: %s", job.Status)
    }
    
    job.Status = "canceled"
    job.UpdatedAt = time.Now()
    
    return store.Update(job)
}
```

## Recurring Jobs

Recurring jobs can be scheduled:

```go
// Schedule a recurring job
func ScheduleRecurringJob(worker *Worker, name string, params interface{}, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        
        for {
            // Enqueue the job
            EnqueueJob(worker.store, name, params, time.Now())
            
            // Wait for next interval
            <-ticker.C
        }
    }()
}

// Example usage
func ScheduleDailyBackup() {
    params := BackupParams{
        Path: "/backups",
    }
    
    ScheduleRecurringJob(worker, "backup.run", params, 24*time.Hour)
}
```
