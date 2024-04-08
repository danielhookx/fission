package fission

import (
	"context"
	"math"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDist struct {
	key int
}

func (d *mockDist) Register(ctx context.Context) {
	return
}
func (d *mockDist) Key() any {
	return d.key
}
func (d *mockDist) Dist(data any) error {
	return nil
}
func (d *mockDist) Close() error {
	return nil
}

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
			u.AddDistributor(&mockDist{key: rand.Intn(math.MaxInt)})
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
	r := NewCenter("test")
	r.AddDistributor(&testDist{
		key: "p1",
		t:   t,
	})
	r.AddDistributor(&testDist{
		key: "p2",
		t:   t,
	})
	r.AddDistributor(&testDist{
		key: "p3",
		t:   t,
	})
	s := "hello world"
	err := r.Fission([]byte(s))
	assert.Nil(t, err)
}
