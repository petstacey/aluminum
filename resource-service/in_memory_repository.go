package main

import "errors"

type InMemoryResourceService struct {
	Workgroups []Workgroup
	JobTitles  []JobTitle
	Locations  []Location
	Resources  []*Resource
}

func NewInMemoryResourceService() *InMemoryResourceService {
	return &InMemoryResourceService{
		Workgroups: []Workgroup{},
		JobTitles:  []JobTitle{},
		Locations:  []Location{},
		Resources:  []*Resource{},
	}
}

func (s *InMemoryResourceService) createResource(resource *Resource) error {
	for _, res := range s.Resources {
		if res.ID == resource.ID {
			return errors.New("duplicate entry")
		}
	}
	s.Resources = append(s.Resources, resource)
	return nil
}

func (s *InMemoryResourceService) getResource(id int64) (*Resource, error) {
	for _, res := range s.Resources {
		if res.ID == id {
			return res, nil
		}
	}
	return nil, errors.New("resource does not exist")
}

func (s *InMemoryResourceService) getResources(name string, titles, types, workgroups, locations, managers []string, filters Filters) ([]*Resource, Metadata, error) {
	return s.Resources, Metadata{}, nil
}

func (s *InMemoryResourceService) updateResource(resource *Resource) error {
	return nil
}
