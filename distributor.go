package fission

import (
	"context"
	"sync"
)

type CreateDistributionHandleFunc func(key any) Distribution

type Distribution interface {
	Register(ctx context.Context)
	Key() any
	Dist(data any) error
	Close() error
}

type DistributorManager struct {
	sync.RWMutex
	distributors map[any]Distribution
}

func NewDistributorManager() *DistributorManager {
	return &DistributorManager{
		distributors: make(map[any]Distribution),
	}
}

func (m *DistributorManager) PutDistributor(key any, distributionCreator CreateDistributionHandleFunc) (p Distribution) {
	var ok bool
	m.Lock()
	if p, ok = m.distributors[key]; !ok && distributionCreator != nil {
		p = distributionCreator(key)
		m.distributors[key] = p
	}
	m.Unlock()
	return
}

func (m *DistributorManager) Destroy() {
	m.Lock()
	for _, platform := range m.distributors {
		platform.Close()
	}
	m.distributors = nil
	m.Unlock()
}
