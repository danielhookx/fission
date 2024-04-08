package fission

import (
	"sync"

	xmap "github.com/danielhookx/xcontainer/map"
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
	distributors *xmap.OrderedMap[any, Distribution]
}

func NewCenter(key any) *Center {
	return &Center{
		key:          key,
		distributors: xmap.NewOrderedMap[any, Distribution](),
	}
}

func (c *Center) AddDistributor(d Distribution) {
	if d == nil {
		return
	}
	c.Lock()
	c.distributors.Set(d.Key(), d)
	c.Unlock()
}

func (c *Center) DelDistributor(key any) {
	c.Lock()
	c.distributors.Delete(key)
	c.Unlock()
}

func (c *Center) Fission(data any) error {
	c.RLock()
	distributors := c.distributors.ToArray()
	c.RUnlock()
	for _, distributor := range distributors {
		err := distributor.Dist(data)
		if err != nil {
			return err
		}
	}
	return nil
}
