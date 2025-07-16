package application

import (
	"context"
	"fmt"
	"hostamat/pkg/platforma/log"
	"sync"
	"time"
)

// ErrStartupTaskFailed represents an error that occurs when a startup task fails.
type ErrStartupTaskFailed struct {
	err error
}

// Error returns the formatted error message for ErrStartupTaskFailed.
func (e *ErrStartupTaskFailed) Error() string {
	return fmt.Sprintf("failed to run startup task: %v", e.err)
}

// Unwrap returns the underlying error for ErrStartupTaskFailed.
func (e *ErrStartupTaskFailed) Unwrap() error {
	return e.err
}

// StartupTaskConfig contains configuration options for a startup task.
type StartupTaskConfig struct {
	Name         string // Name of the startup task
	AbortOnError bool   // Whether to abort application startup if this task fails
}

// startupTask represents an individual startup task with its runner and configuration.
type startupTask struct {
	runner Runner
	config StartupTaskConfig
}

// Application manages startup tasks and services for the application lifecycle.
type Application struct {
	startupTasks []startupTask
	services     map[string]Runner
	health       *ApplicationHealth
}

// New creates and returns a new Application instance.
func New() *Application {
	return &Application{services: make(map[string]Runner), health: NewApplicationHealth()}
}

func (a *Application) Health() *ApplicationHealth {
	return a.health
}

// OnStart registers a new startup task with the given runner and configuration.
func (a *Application) OnStart(task Runner, config StartupTaskConfig) {
	a.startupTasks = append(a.startupTasks, startupTask{task, config})
}

// Register adds a named service to the application.
func (a *Application) Register(serviceName string, service Runner) {
	a.services[serviceName] = service
	a.health.Services[serviceName] = &ServiceHealth{Status: ServiceStatusNotStarted}
}

// Run executes all startup tasks and services in the application.
// It returns an error if any startup task configured to abort on error fails.
func (a *Application) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	log.InfoContext(ctx, "Starting application", "startupTasks", len(a.startupTasks))

	for i, task := range a.startupTasks {
		log.InfoContext(ctx, "running task", "task", task.config.Name, "index", i)

		taskCtx := context.WithValue(ctx, log.SartupTaskKey, task.config.Name)

		err := task.runner.Run(taskCtx)
		if err != nil {
			log.ErrorContext(ctx, "error in startup task", "error", err, "task", task.config.Name)

			if task.config.AbortOnError {
				return &ErrStartupTaskFailed{err: err}
			}
		}
	}

	var wg sync.WaitGroup

	for serviceName, service := range a.services {
		wg.Add(1)

		serviceCtx := context.WithValue(ctx, log.ServiceNameKey, serviceName)

		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.ErrorContext(serviceCtx, "service panicked", "service", serviceName, "panic", r)
				}
			}()

			log.InfoContext(ctx, "starting service", "service", serviceName)
			a.health.StartService(serviceName)
			log.DebugContext(ctx, "health", "health", a.health)

			err := service.Run(serviceCtx)
			if err != nil {
				a.health.FailService(serviceName, err)
				log.ErrorContext(ctx, "error in service", "service", serviceName, "error", err)
			}
		}()
	}

	a.health.StartedAt = time.Now()

	wg.Wait()

	return nil
}
