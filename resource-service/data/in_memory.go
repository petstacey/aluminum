package data

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

func (s *InMemoryResourceService) CreateResource(resource *Resource) error {
	for _, res := range s.Resources {
		if res.ID == resource.ID {
			return errors.New("duplicate entry")
		}
	}
	s.Resources = append(s.Resources, resource)
	return nil
}

func (s *InMemoryResourceService) GetResource(id int64) (*Resource, error) {
	for _, res := range s.Resources {
		if res.ID == id {
			return res, nil
		}
	}
	return nil, errors.New("resource does not exist")
}

func (s *InMemoryResourceService) GetResources(name string, titles, types, workgroups, locations, managers []string, filters Filters) ([]*Resource, Metadata, error) {
	return s.Resources, Metadata{}, nil
}

func (s *InMemoryResourceService) UpdateResource(resource *Resource) error {
	return nil
}
