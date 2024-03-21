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

func (tc *testDistCreator) Create(ctx context.Context, pid any) Distribution {
	ep := ctx.Value("endpoint")
	return &testDist{
		id: pid.(string),
		ep: ep.(string),
		t:  tc.t,
	}
}

type testDist struct {
	id string
	ep string
	t  testing.TB
}

func (td *testDist) Dist(data any) error {
	td.t.Logf("%s--received:%s, publish to %s\n", td.id, string(data.([]byte)), td.ep)
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
			pm.PutDistributor(String(uuid.NewString()), "", tc.Create)
		}
	})
}

func createDoNothingDist(ctx context.Context, pid any) Distribution {
	return &doNothingDist{}
}

type doNothingDist struct{}

func (d *doNothingDist) Dist(data any) error { return nil }
func (d *doNothingDist) Close() error        { return nil }

func BenchmarkPlatform_Push(b *testing.B) {
	pm := NewDistributorManager()
	p := pm.PutDistributor(String("p1"), "e1", createDoNothingDist)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Push([]byte("hello"))
		}
	})
}

func TestPlatform_Push(t *testing.T) {
	tc := testDistCreator{t: t}
	pm := NewDistributorManager()
	p := pm.PutDistributor(String("p1"), "e1", tc.Create)
	err := p.Push([]byte("hello platform p1"))
	assert.Nil(t, err)
}
