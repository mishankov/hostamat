package scheduler

import (
	"context"
	"hostamat/pkg/platforma/log"
	"time"

	"github.com/google/uuid"
)

// SchedulerAction defines the function signature for actions executed by the scheduler.
type SchedulerAction func(context.Context) error

// Scheduler represents a periodic task runner that executes an action at fixed intervals.
type Scheduler struct {
	period time.Duration   // The interval between action executions
	action SchedulerAction // The function to execute periodically
}

// New creates a new Scheduler instance with the specified period and action.
func New(period time.Duration, action SchedulerAction) *Scheduler {
	return &Scheduler{period: period, action: action}
}

// Run starts the scheduler and executes the action at the configured interval.
// The scheduler will continue running until the context is canceled.
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runCtx := context.WithValue(ctx, log.TraceIdKey, uuid.NewString())
			err := s.action(runCtx)
			if err != nil {
				log.ErrorContext(ctx, "error in scheduler", "error", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
