package application

import (
	"context"
	"log"
	"sync"
)

type ApplicationService interface {
	Run(context.Context) error
}

type Application struct {
	services map[string]ApplicationService
}

func New() *Application {
	return &Application{services: make(map[string]ApplicationService)}
}

func (a *Application) Register(serviceName string, service ApplicationService) {
	a.services[serviceName] = service
}

func (a *Application) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	if ctx == nil {
		ctx = context.Background()
	}

	for serviceName, service := range a.services {
		wg.Add(1)
		go func() {
			log.Println("starting service", serviceName)
			err := service.Run(ctx)
			if err != nil {
				log.Println("Error in service:", serviceName, err)
			}
		}()
	}

	wg.Wait()

	return nil
}
