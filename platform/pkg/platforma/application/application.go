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
	services []ApplicationService
}

func New() *Application {
	return &Application{}
}

func (a *Application) Register(service ApplicationService) {
	a.services = append(a.services, service)
}

func (a *Application) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	if ctx == nil {
		ctx = context.Background()
	}

	for _, service := range a.services {
		wg.Add(1)
		go func() {
			err := service.Run(ctx)
			if err != nil {
				log.Println("Error in service:", err)
			}
		}()
	}

	wg.Wait()

	return nil
}
