package petrock_example_feature_name

import (
	"context"
	"log/slog"
	"time"

	"petrock_example_module_path/core" // Placeholder for target project's core package
)

// Jobs holds dependencies for background tasks related to the feature.
type Jobs struct {
	state *State           // Example: Access to feature state
	log   *core.MessageLog // Example: Ability to append new commands/events
	// Add other dependencies like external clients, config, etc.
}

// NewJobs creates a new Jobs instance.
func NewJobs(state *State, log *core.MessageLog) *Jobs {
	return &Jobs{
		state: state,
		log:   log,
	}
}

// StartExampleWorker demonstrates a long-running background worker.
// This function would typically be launched as a goroutine from the application's main process (e.g., cmd/serve.go).
func (j *Jobs) StartExampleWorker(ctx context.Context) {
	slog.Info("Starting example worker", "feature", "petrock_example_feature_name")
	ticker := time.NewTicker(1 * time.Minute) // Example: Run every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			slog.Debug("Example worker tick", "feature", "petrock_example_feature_name")
			// --- Worker Logic ---
			// Example: Find items needing processing based on state
			// itemsToProcess := j.findItemsNeedingWork()
			// for _, item := range itemsToProcess {
			//     err := j.processItem(ctx, item)
			//     if err != nil {
			//         slog.Error("Failed to process item in worker", "error", err, "itemID", item.ID)
			//     }
			// }
			// --- End Worker Logic ---

		case <-ctx.Done():
			slog.Info("Stopping example worker due to context cancellation", "feature", "petrock_example_feature_name")
			return
		}
	}
}

// StartExampleScheduledTask demonstrates a task running on a schedule.
// This might use a cron library or simple time checks.
func (j *Jobs) StartExampleScheduledTask(ctx context.Context /*, schedule string */) {
	slog.Info("Starting example scheduled task", "feature", "petrock_example_feature_name")
	// Example using a simple ticker (replace with cron library for complex schedules)
	ticker := time.NewTicker(24 * time.Hour) // Example: Run daily
	defer ticker.Stop()

	// Run once immediately on start?
	j.runScheduledTaskLogic(ctx)

	for {
		select {
		case <-ticker.C:
			j.runScheduledTaskLogic(ctx)
		case <-ctx.Done():
			slog.Info("Stopping example scheduled task due to context cancellation", "feature", "petrock_example_feature_name")
			return
		}
	}
}

// runScheduledTaskLogic contains the actual work performed by the scheduled task.
func (j *Jobs) runScheduledTaskLogic(ctx context.Context) {
	slog.Debug("Running scheduled task logic", "feature", "petrock_example_feature_name")
	// --- Task Logic ---
	// Example: Generate a daily report based on feature state
	// reportData, err := j.generateReport(ctx)
	// if err != nil {
	//     slog.Error("Failed to generate report", "error", err)
	//     return
	// }
	// err = j.sendReport(ctx, reportData)
	// if err != nil {
	//     slog.Error("Failed to send report", "error", err)
	// }
	// --- End Task Logic ---
}

// Add helper functions for worker/task logic below...
// func (j *Jobs) findItemsNeedingWork() []*Item { ... }
// func (j *Jobs) processItem(ctx context.Context, item *Item) error { ... }
// func (j *Jobs) generateReport(ctx context.Context) (string, error) { ... }
// func (j *Jobs) sendReport(ctx context.Context, report string) error { ... }
