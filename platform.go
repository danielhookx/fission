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

type PlatformManager struct {
	sync.RWMutex
	platforms           map[any]*Platform
	distributionCreator CreateDistributionHandleFunc
}

func NewPlatformManager(distributionCreator CreateDistributionHandleFunc) *PlatformManager {
	return &PlatformManager{
		platforms:           make(map[any]*Platform),
		distributionCreator: distributionCreator,
	}
}

func (m *PlatformManager) PutPlatform(ctx context.Context, key any, forceDC CreateDistributionHandleFunc) (p *Platform) {
	var ok bool
	dc := m.distributionCreator
	if forceDC != nil {
		dc = forceDC
	}
	m.Lock()
	if p, ok = m.platforms[key]; !ok {
		p = NewPlatform(ctx, key, dc)
		m.platforms[key] = p
	}
	m.Unlock()
	return
}

func (m *PlatformManager) Destroy() {
	m.Lock()
	for _, platform := range m.platforms {
		platform.Close()
	}
	m.platforms = nil
	m.Unlock()
}

type Platform struct {
	key          any
	distribution Distribution
}

func NewPlatform(ctx context.Context, key any, createDist CreateDistributionHandleFunc) *Platform {
	dist := createDist(ctx, key)
	return &Platform{
		key:          key,
		distribution: dist,
	}
}

// Push push data to the dispatcher, the dispatcher handles the data distribution.
func (p *Platform) Push(data any) error {
	return p.distribution.Dist(data)
}

func (p *Platform) Close() error {
	return p.distribution.Close()
}
