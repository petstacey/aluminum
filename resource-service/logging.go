package main

import (
	"fmt"
	"time"

	"github.com/petstacey/aluminum/resource-service/data"
)

type LoggingService struct {
	next Service
}

// createResource implements Service.
func (s *LoggingService) CreateResource(resource *data.Resource) (err error) {
	defer func(start time.Time) {
		fmt.Printf("create resource: %v, err: %v, took: %v\n", resource, err, time.Since(start))
	}(time.Now())
	return s.next.CreateResource(resource)
}

// getResource implements Service.
func (s *LoggingService) GetResource(id int64) (res *data.Resource, err error) {
	defer func(start time.Time) {
		fmt.Printf("get resource with id: %v, err: %v, took: %v\n", id, err, time.Since(start))
	}(time.Now())
	return s.next.GetResource(id)
}

// getResources implements Service.
func (s *LoggingService) GetResources(name string, titles, types, workgroups, locations, managers []string, filters data.Filters) (res []*data.Resource, metadata data.Metadata, err error) {
	defer func(start time.Time) {
		fmt.Printf("get resources: %v, err: %v, took: %v\n", res, err, time.Since(start))
	}(time.Now())
	return s.next.GetResources(name, titles, types, workgroups, locations, managers, filters)
}

// updateResource implements Service.
func (s *LoggingService) UpdateResource(resource *data.Resource) (err error) {
	defer func(start time.Time) {
		fmt.Printf("update resources: %v, err: %s, took: %v\n", resource, err.Error(), time.Since(start))
	}(time.Now())
	return s.next.UpdateResource(resource)
}

func NewLoggingService(next Service) Service {
	return &LoggingService{
		next: next,
	}
}
