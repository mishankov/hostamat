package application

import (
	"encoding/json"
	"time"
)

type ServiceStatus string

const (
	ServiceStatusNotStarted ServiceStatus = "NOT_STARTED"
	ServiceStatusStarted    ServiceStatus = "STARTED"
	ServiceStatusError      ServiceStatus = "ERROR"
)

type ServiceHealth struct {
	Status    ServiceStatus `json:"status"`
	StartedAt time.Time     `json:"startedAt"`
	StoppedAt time.Time     `json:"stoppedAt,omitempty"`
	Error     string        `json:"error,omitempty"`
}

type ApplicationHealth struct {
	StartedAt time.Time                 `json:"startedAt"`
	Services  map[string]*ServiceHealth `json:"services"`
}

func NewApplicationHealth() *ApplicationHealth {
	return &ApplicationHealth{Services: make(map[string]*ServiceHealth)}
}

func (h *ApplicationHealth) StartService(serviceName string) {
	if service, ok := h.Services[serviceName]; ok {
		service.Status = ServiceStatusStarted
		service.StartedAt = time.Now()

		h.Services[serviceName] = service
	}
}

func (h *ApplicationHealth) FailService(serviceName string, err error) {
	if service, ok := h.Services[serviceName]; ok {
		service.Status = ServiceStatusError
		service.StoppedAt = time.Now()
		service.Error = err.Error()

		h.Services[serviceName] = service
	}
}

func (h *ApplicationHealth) String() string {
	b, _ := json.Marshal(h)
	return string(b)
}

func (h *ApplicationHealth) StartApplication() {
	h.StartedAt = time.Now()
}
