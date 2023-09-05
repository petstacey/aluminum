package main

type Service interface {
	getResource(id int64) (*Resource, error)
	getResources(name string, titles, types, workgroups, locations, managers []string, filters Filters) ([]*Resource, Metadata, error)
	createResource(resource *Resource) error
	updateResource(resource *Resource) error
}

type ResourceService struct {
	repo Service
}

func NewService(repo Service) Service {
	return &ResourceService{
		repo: repo,
	}
}

func (s *ResourceService) createResource(resource *Resource) error {
	return s.repo.createResource(resource)
}

func (s *ResourceService) getResource(id int64) (*Resource, error) {
	return s.repo.getResource(id)
}

func (s *ResourceService) getResources(name string, titles, types, workgroups, locations, managers []string, filters Filters) ([]*Resource, Metadata, error) {
	return s.repo.getResources(name, titles, types, workgroups, locations, managers, filters)
}

func (s *ResourceService) updateResource(resource *Resource) error {
	return s.repo.updateResource(resource)
}
