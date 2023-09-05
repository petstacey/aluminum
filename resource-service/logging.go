package main

import (
	"fmt"
	"time"
)

type LoggingService struct {
	next Service
}

// createResource implements Service.
func (s *LoggingService) createResource(resource *Resource) (err error) {
	defer func(start time.Time) {
		fmt.Printf("create resource: %v, err: %v, took: %v\n", resource, err, time.Since(start))
	}(time.Now())
	return s.next.createResource(resource)
}

// getResource implements Service.
func (s *LoggingService) getResource(id int64) (res *Resource, err error) {
	defer func(start time.Time) {
		fmt.Printf("get resource with id: %v, err: %v, took: %v\n", id, err, time.Since(start))
	}(time.Now())
	return s.next.getResource(id)
}

// getResources implements Service.
func (s *LoggingService) getResources(name string, titles, types, workgroups, locations, managers []string, filters Filters) (res []*Resource, metadata Metadata, err error) {
	defer func(start time.Time) {
		fmt.Printf("get resources: %v, err: %v, took: %v\n", res, err, time.Since(start))
	}(time.Now())
	return s.next.getResources(name, titles, types, workgroups, locations, managers, filters)
}

// updateResource implements Service.
func (s *LoggingService) updateResource(resource *Resource) (err error) {
	defer func(start time.Time) {
		fmt.Printf("update resources: %v, err: %s, took: %v\n", resource, err.Error(), time.Since(start))
	}(time.Now())
	return s.next.updateResource(resource)
}

func NewLoggingService(next Service) Service {
	return &LoggingService{
		next: next,
	}
}
