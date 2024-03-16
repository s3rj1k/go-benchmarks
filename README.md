# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/local/char-xxhash-matrix-16         	    6742	    174538 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16             	     853	   1443185 ns/op	  276010 B/op	     446 allocs/op
BenchmarkSets/ironpark/skiplist-16                	    6714	    169718 ns/op	   25694 B/op	     765 allocs/op
BenchmarkSets/alphadose/haxmap-16                 	   23595	     48414 ns/op	   12240 B/op	     255 allocs/op
BenchmarkSets/dolthub/swiss-16                    	   92416	     12837 ns/op	       0 B/op	       0 allocs/op
BenchmarkSets/runtime/map-16                      	   55465	     21480 ns/op	      76 B/op	       0 allocs/op
```
