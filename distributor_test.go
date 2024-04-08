package fission

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func String(v string) context.Context {
	return context.WithValue(context.Background(), "endpoint", v)
}

type testDistCreator struct {
	t testing.TB
}

func (tc *testDistCreator) Create(key any) Distribution {
	return &testDist{
		key: key.(string),
		t:   tc.t,
	}
}

type testDist struct {
	key string
	t   testing.TB
}

func (td *testDist) Register(ctx context.Context) {
	return
}
func (td *testDist) Key() any {
	return td.key
}

func (td *testDist) Dist(data any) error {
	td.t.Logf("%s--received:%s\n", td.key, string(data.([]byte)))
	return nil
}

func (td *testDist) Close() error {
	return nil
}

func BenchmarkPlatformManager_PutPlatform(b *testing.B) {
	tc := testDistCreator{t: b}
	pm := NewDistributorManager()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pm.PutDistributor(uuid.NewString(), tc.Create)
		}
	})
}

func createDoNothingDist(pid any) Distribution {
	return &doNothingDist{}
}

type doNothingDist struct {
	key int
}

func (d *doNothingDist) Register(ctx context.Context) { return }
func (d *doNothingDist) Key() any                     { return d.key }
func (d *doNothingDist) Dist(data any) error          { return nil }
func (d *doNothingDist) Close() error                 { return nil }

func BenchmarkPlatform_Push(b *testing.B) {
	pm := NewDistributorManager()
	p := pm.PutDistributor("p1", createDoNothingDist)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Dist([]byte("hello"))
		}
	})
}

func TestPlatform_Push(t *testing.T) {
	tc := testDistCreator{t: t}
	pm := NewDistributorManager()
	p := pm.PutDistributor("p1", tc.Create)
	err := p.Dist([]byte("hello platform p1"))
	assert.Nil(t, err)
}
