package application

import (
	"context"
	"hostamat/pkg/platforma/log"
	"sync"
)

type Runner interface {
	Run(context.Context) error
}

type Application struct {
	services map[string]Runner
}

func New() *Application {
	return &Application{services: make(map[string]Runner)}
}

func (a *Application) Register(serviceName string, service Runner) {
	a.services[serviceName] = service
}

func (a *Application) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	if ctx == nil {
		ctx = context.Background()
	}

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
