package main

type Service struct {
	Name   string
	Checks []HealthCheck
	Tags   []string
}

type HealthCheck struct {
	Status  int
	Message int
}
