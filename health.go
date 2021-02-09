package main

import (
	"fmt"
)

type Service struct {
	Name   string        `json:"name"`
	Checks []HealthCheck `json:"health_checks"`
	Status int           `json:"status"` // 0 - Operational, 1 - Degradation, 2 - Downtime, 3 - Unknown
	Tags   []string      `json:"tags"`
	Feed   []Outage      `json:"feed"`
}

func (service Service) GetDescription() (string, []string) {
	return service.Name, service.Tags
}

func (service Service) GetShortService() Service {
	service.Checks = []HealthCheck{}
	return service
}

type HealthCheck struct {
	Status       int           `json:"status"` // 0 - OK, 1 - WARN, 2 - CRITICAL, 3 - UNKNOWN
	Title        string        `json:"title"`
	Details      string        `json:"details"`
	Priority     int           `json:"priority"`
	Relationship []HealthCheck `json:"rels"`
}

func (check *HealthCheck) SetStatus(status int) error {
	possibleValues := []int{0, 1, 2, 3}
	for _, value := range possibleValues {
		if status == value {
			check.Status = status
			return nil
		}
	}
	err := fmt.Errorf("Status %v is wrong, possible values are", possibleValues)
	return err
}

type Outage struct {
	Summary string `json:"summary"`
	Details string `json:"details"`
	Start   string `json:"scheduled_begin"`
	End     string `json:"schedule_end"`
}
