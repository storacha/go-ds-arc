package arc

import (
	"testing"

	dstest "github.com/ipfs/go-datastore/test"
)

func TestARCDS(t *testing.T) {
	ds := New(400)
	dstest.SubtestAll(t, ds)
}
