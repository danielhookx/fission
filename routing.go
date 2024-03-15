package fission

import (
	"sync"
)

type RouteManager struct {
	sync.RWMutex
	routes map[any]*Route
}

func NewRouteManager() *RouteManager {
	return &RouteManager{
		routes: make(map[any]*Route),
	}
}

func (m *RouteManager) PutRoute(key any) (r *Route) {
	var ok bool
	m.Lock()
	if r, ok = m.routes[key]; !ok {
		r = NewRoute(key)
		m.routes[key] = r
	}
	m.Unlock()
	return
}

func (m *RouteManager) Destroy() {
	m.Lock()
	m.routes = nil
	m.Unlock()
}

type Route struct {
	sync.RWMutex
	key       any
	platforms map[any]*Platform
}

func NewRoute(key any) *Route {
	return &Route{
		key:       key,
		platforms: make(map[any]*Platform),
	}
}

func (r *Route) AddPlatform(p *Platform) {
	if p == nil {
		return
	}
	r.Lock()
	r.platforms[p.key] = p
	r.Unlock()
}

func (r *Route) DelPlatform(key any) {
	r.Lock()
	delete(r.platforms, key)
	r.Unlock()
}

func (r *Route) GetPlatforms() map[any]*Platform {
	r.RLock()
	defer r.RUnlock()
	return r.platforms
}

func (r *Route) Fission(data any) error {
	r.RLock()
	defer r.RUnlock()
	for _, platform := range r.platforms {
		err := platform.Push(data)
		if err != nil {
			return err
		}
	}
	return nil
}
