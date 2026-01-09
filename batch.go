package arc

import (
	"context"

	"github.com/ipfs/go-datastore"
)

type op struct {
	delete bool
	value  []byte
}

type batch struct {
	ds  *ARCDS
	ops map[datastore.Key]op
}

var _ datastore.Batch = (*batch)(nil)

func (bt *batch) Put(ctx context.Context, key datastore.Key, val []byte) error {
	bt.ops[key] = op{value: val}
	return nil
}

func (bt *batch) Delete(ctx context.Context, key datastore.Key) error {
	bt.ops[key] = op{delete: true}
	return nil
}

func (bt *batch) Commit(ctx context.Context) error {
	bt.ds.lock.Lock()
	for k, op := range bt.ops {
		if op.delete {
			bt.ds.arc.Remove(k.String())
		} else {
			bt.ds.arc.Add(k.String(), op.value)
		}
	}
	bt.ds.lock.Unlock()
	return nil
}
