# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/local/char-xxhash-matrix-16         	    6477	    182545 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16             	     819	   1442124 ns/op	  279886 B/op	     454 allocs/op
BenchmarkSets/ironpark/skiplist-16                	    7112	    167491 ns/op	   25694 B/op	     765 allocs/op
BenchmarkSets/alphadose/haxmap-16                 	   24506	     48082 ns/op	   12240 B/op	     255 allocs/op
BenchmarkSets/runtime/map-16                      	   55561	     21660 ns/op	      79 B/op	       0 allocs/op
```
