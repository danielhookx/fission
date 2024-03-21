package fission

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func BenchmarkRouteManager_PutRoute(b *testing.B) {
	rm := NewCenterManager()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rm.PutCenter(uuid.NewString())
		}
	})
}

func BenchmarkRoute_AddPlatform(b *testing.B) {
	rm := NewCenterManager()
	u := rm.PutCenter("1")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u.AddDistributor(&Distributor{key: uuid.NewString()})
		}
	})
}

func BenchmarkRoute_DelPlatform(b *testing.B) {
	rm := NewCenterManager()
	u := rm.PutCenter("1")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u.DelDistributor(uuid.NewString())
		}
	})
}

func TestRoute_Fission(t *testing.T) {
	r := Center{
		key: "test",
		distributors: map[any]*Distributor{
			"p1": {
				key: "p1",
				distribution: &testDist{
					id: "p1",
					t:  t,
				},
			},
			"p2": {
				key: "p2",
				distribution: &testDist{
					id: "p2",
					t:  t,
				},
			},
			"p3": {
				key: "p3",
				distribution: &testDist{
					id: "p3",
					t:  t,
				},
			},
		},
	}
	s := "hello world"
	err := r.Fission([]byte(s))
	assert.Nil(t, err)
}
