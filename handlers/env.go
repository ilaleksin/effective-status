package handlers

import "effective-status/model"

type Env struct {
	services
	deps
}

type services interface {
	All() ([]model.Service, error)
	Create(model.Service) (int, error)
	Update(string, model.Service) (bool, error)
	Get(string) (model.Service, error)
	Delete(string) (bool, error)
}

type deps interface {
	Get(string) ([]model.Dependency, error)
	Create(model.Dependency) (int, error)
	Delete(int) (bool, error)
}
