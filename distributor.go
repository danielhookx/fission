package fission

import (
	"context"
	"sync"
)

type CreateDistributionHandleFunc func(ctx context.Context, key any) Distribution

type Distribution interface {
	Dist(data any) error
	Close() error
}

type DistributorManager struct {
	sync.RWMutex
	distributors map[any]*Distributor
}

func NewDistributorManager() *DistributorManager {
	return &DistributorManager{
		distributors: make(map[any]*Distributor),
	}
}

func (m *DistributorManager) PutDistributor(ctx context.Context, key any, distributionCreator CreateDistributionHandleFunc) (p *Distributor) {
	var ok bool
	m.Lock()
	if p, ok = m.distributors[key]; !ok {
		p = NewDistributor(ctx, key, distributionCreator)
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

type Distributor struct {
	key          any
	distribution Distribution
}

func NewDistributor(ctx context.Context, key any, createDist CreateDistributionHandleFunc) *Distributor {
	dist := createDist(ctx, key)
	return &Distributor{
		key:          key,
		distribution: dist,
	}
}

// Push push data to the dispatcher, the dispatcher handles the data distribution.
func (p *Distributor) Push(data any) error {
	return p.distribution.Dist(data)
}

func (p *Distributor) Close() error {
	return p.distribution.Close()
}
