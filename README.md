# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/local/char-xxhash-matrix-16         	    6831	    173816 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16             	     895	   1350956 ns/op	  269827 B/op	     448 allocs/op
BenchmarkSets/ironpark/skiplist-16                	    6825	    169648 ns/op	   25695 B/op	     765 allocs/op
BenchmarkSets/runtime/map-16                      	   55693	     21905 ns/op	      78 B/op	       0 allocs/op
```
