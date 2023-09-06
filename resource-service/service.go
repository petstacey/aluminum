package main

import "github.com/petstacey/aluminum/resource-service/data"

type Service interface {
	GetResource(id int64) (*data.Resource, error)
	GetResources(name string, titles, types, workgroups, locations, managers []string, filters data.Filters) ([]*data.Resource, data.Metadata, error)
	CreateResource(resource *data.Resource) error
	UpdateResource(resource *data.Resource) error
}

type ResourceService struct {
	repo Service
}

func NewService(repo Service) Service {
	return &ResourceService{
		repo: repo,
	}
}

func (s *ResourceService) CreateResource(resource *data.Resource) error {
	return s.repo.CreateResource(resource)
}

func (s *ResourceService) GetResource(id int64) (*data.Resource, error) {
	return s.repo.GetResource(id)
}

func (s *ResourceService) GetResources(name string, titles, types, workgroups, locations, managers []string, filters data.Filters) ([]*data.Resource, data.Metadata, error) {
	return s.repo.GetResources(name, titles, types, workgroups, locations, managers, filters)
}

func (s *ResourceService) UpdateResource(resource *data.Resource) error {
	return s.repo.UpdateResource(resource)
}
