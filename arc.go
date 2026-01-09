package arc

import (
	"context"
	"sync"

	"github.com/hashicorp/golang-lru/arc/v2"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

type ARCDS struct {
	arc  *arc.ARCCache[string, []byte]
	lock sync.RWMutex
}

func (a *ARCDS) Batch(ctx context.Context) (datastore.Batch, error) {
	return &batch{ds: a, ops: map[datastore.Key]op{}}, nil
}

func (a *ARCDS) Close() error {
	return nil
}

func (a *ARCDS) Delete(ctx context.Context, key datastore.Key) error {
	a.lock.Lock()
	a.arc.Remove(key.String())
	a.lock.Unlock()
	return nil
}

func (a *ARCDS) Get(ctx context.Context, key datastore.Key) (value []byte, err error) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	v, ok := a.arc.Get(key.String())
	if !ok {
		return nil, datastore.ErrNotFound
	}
	return v, nil
}

func (a *ARCDS) GetSize(ctx context.Context, key datastore.Key) (size int, err error) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	v, ok := a.arc.Peek(key.String())
	if !ok {
		return -1, datastore.ErrNotFound
	}
	return len(v), nil
}

func (a *ARCDS) Has(ctx context.Context, key datastore.Key) (exists bool, err error) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.arc.Contains(key.String()), nil
}

func (a *ARCDS) Put(ctx context.Context, key datastore.Key, value []byte) error {
	a.lock.Lock()
	a.arc.Add(key.String(), value)
	a.lock.Unlock()
	return nil
}

func (a *ARCDS) Query(ctx context.Context, q query.Query) (query.Results, error) {
	a.lock.RLock()
	keys := a.arc.Keys()
	i := 0
	var once sync.Once
	it := query.Iterator{
		Next: func() (query.Result, bool) {
			if i >= len(keys) {
				once.Do(a.lock.RUnlock)
				return query.Result{}, false
			}
			k := keys[i]
			i++
			val, _ := a.arc.Get(k)
			ent := query.Entry{Key: k, Value: val, Size: len(val)}
			return query.Result{Entry: ent}, true
		},
		Close: func() error {
			once.Do(a.lock.RUnlock)
			return nil
		},
	}
	return query.NaiveQueryApply(q, query.ResultsFromIterator(q, it)), nil
}

func (a *ARCDS) Sync(ctx context.Context, prefix datastore.Key) error {
	return nil
}

func New(size int) *ARCDS {
	arcCache, _ := arc.NewARC[string, []byte](size)
	return &ARCDS{arc: arcCache}
}

var _ datastore.Batching = (*ARCDS)(nil)
