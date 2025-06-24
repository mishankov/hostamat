package application

import (
	"context"
	"hostamat/pkg/platforma/log"
	"sync"
)

type Runner interface {
	Run(context.Context) error
}

type RunnerFunc func(context.Context) error

func (rf RunnerFunc) Run(ctx context.Context) error {
	return rf(ctx)
}

type StartupTaskConfig struct {
	Name         string
	AbortOnError bool
}

type startupTask struct {
	runner Runner
	config StartupTaskConfig
}

type Application struct {
	startupTasks []startupTask
	services     map[string]Runner
}

func New() *Application {
	return &Application{services: make(map[string]Runner)}
}

func (a *Application) OnStart(task Runner, config StartupTaskConfig) {
	a.startupTasks = append(a.startupTasks, startupTask{task, config})
}

func (a *Application) Register(serviceName string, service Runner) {
	a.services[serviceName] = service
}

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
				// TODO: wrap error
				return err
			}
		}
	}

	var wg sync.WaitGroup

	for serviceName, service := range a.services {
		wg.Add(1)
		serviceCtx := context.WithValue(ctx, log.ServiceNameKey, serviceName)
		go func() {
			log.InfoContext(ctx, "starting service", "service", serviceName)
			err := service.Run(serviceCtx)
			if err != nil {
				log.ErrorContext(ctx, "error in service", "service", serviceName, "error", err)
			}
		}()
	}

	wg.Wait()

	return nil
}
