package scheduler

import (
	"context"
	"hostamat/pkg/platforma/application"
	"hostamat/pkg/platforma/log"
	"time"

	"github.com/google/uuid"
)

type logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

// SchedulerAction defines the function signature for actions executed by the scheduler.
type SchedulerAction func(context.Context) error

// Scheduler represents a periodic task runner that executes an action at fixed intervals.
type Scheduler struct {
	logger logger             // The logger for logging messages
	period time.Duration      // The interval between action executions
	runner application.Runner // The runner to execute periodically
}

// New creates a new Scheduler instance with the specified period and action.
func New(logger logger, period time.Duration, runner application.Runner) *Scheduler {
	return &Scheduler{logger: logger, period: period, runner: runner}
}

// Run starts the scheduler and executes the runner at the configured interval.
// The scheduler will continue running until the context is canceled.
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.InfoContext(ctx, "scheduler task started")
			runCtx := context.WithValue(ctx, log.TraceIdKey, uuid.NewString())
			err := s.runner.Run(runCtx)
			if err != nil {
				s.logger.ErrorContext(ctx, "error in scheduler", "error", err)
			}
			s.logger.InfoContext(ctx, "scheduler task finished")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
