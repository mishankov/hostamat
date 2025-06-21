package scheduler

import (
	"context"
	"hostamat/pkg/platforma/log"
	"time"

	"github.com/google/uuid"
)

type SchedulerAction func(context.Context) error

type Scheduler struct {
	period time.Duration
	action SchedulerAction
}

func New(period time.Duration, action SchedulerAction) *Scheduler {
	return &Scheduler{period: period, action: action}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runCtx := context.WithValue(ctx, "traceId", uuid.NewString())
			err := s.action(runCtx)
			if err != nil {
				log.ErrorContext(ctx, "error in scheduler", "error", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
