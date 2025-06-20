package application

import (
	"log"
	"sync"
)

type ApplicationService interface {
	Run() error
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

func (a *Application) Run() error {
	var wg sync.WaitGroup

	for _, service := range a.services {
		wg.Add(1)
		go func() {
			err := service.Run()
			if err != nil {
				log.Println("Error in service:", err)
			}
		}()
	}

	wg.Wait()

	return nil
}
