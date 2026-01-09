# go-ds-arc

This is an IPFS datastore implementation backed by an Adaptive Replacement
Cache (ARC).

From [golang-lru/arc]():

> ARC is an enhancement over the standard LRU cache in that tracks both
> frequency and recency of use. This avoids a burst in access to new
> entries from evicting the frequently used older entries. It adds some
> additional tracking overhead to a standard LRU cache, computationally
> it is roughly 2x the cost, and the extra memory overhead is linear
> with the size of the cache. ARC has been patented by IBM, but is
> similar to the TwoQueueCache (2Q) which requires setting parameters.

## Install

```sh
go get github.com/storacha/go-ds-arc
```

## Usage

```go
package main

import "github.com/storacha/go-ds-arc"

func main() {
  capacity := 1000

  // create a datastore with capacity for 1,000 items
  ds, err := arc.New(capacity)

  // use as per https://github.com/ipfs/go-datastore
}
```

### API

[pkg.go.dev Reference](https://pkg.go.dev/github.com/storacha/go-ds-arc)

## Contributing

Feel free to join in. All welcome. Please
[open an issue](https://github.com/storacha/go-ds-arc/issues)!

## License

Dual-licensed under [MIT OR Apache 2.0](LICENSE.md)
