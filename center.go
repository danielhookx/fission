package fission

import (
	"sync"
)

type CenterManager struct {
	sync.RWMutex
	centers map[any]*Center
}

func NewCenterManager() *CenterManager {
	return &CenterManager{
		centers: make(map[any]*Center),
	}
}

func (m *CenterManager) PutCenter(key any) (c *Center) {
	var ok bool
	m.Lock()
	if c, ok = m.centers[key]; !ok {
		c = NewCenter(key)
		m.centers[key] = c
	}
	m.Unlock()
	return
}

func (m *CenterManager) Destroy() {
	m.Lock()
	m.centers = nil
	m.Unlock()
}

type Center struct {
	sync.RWMutex
	key          any
	distributors map[any]*Distributor
}

func NewCenter(key any) *Center {
	return &Center{
		key:          key,
		distributors: make(map[any]*Distributor),
	}
}

func (c *Center) AddDistributor(d *Distributor) {
	if d == nil {
		return
	}
	c.Lock()
	c.distributors[d.key] = d
	c.Unlock()
}

func (c *Center) DelDistributor(key any) {
	c.Lock()
	delete(c.distributors, key)
	c.Unlock()
}

func (c *Center) GetDistributors() map[any]*Distributor {
	c.RLock()
	defer c.RUnlock()
	return c.distributors
}

func (c *Center) Fission(data any) error {
	c.RLock()
	defer c.RUnlock()
	for _, distributor := range c.distributors {
		err := distributor.Push(data)
		if err != nil {
			return err
		}
	}
	return nil
}
