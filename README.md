# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/Workiva/go-datastructures/trie/ctrie-16         	    4220	    291564 ns/op	  269465 B/op	    4914 allocs/op
BenchmarkSets/local/char-xxhash-matrix-16                     	    5953	    201552 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16                         	     988	   1135091 ns/op	  288643 B/op	     454 allocs/op
BenchmarkSets/ironpark/skiplist-16                            	    6962	    164425 ns/op	   25695 B/op	     765 allocs/op
BenchmarkSets/alphadose/haxmap-16                             	   26445	     45614 ns/op	   12240 B/op	     255 allocs/op
BenchmarkSets/dolthub/swiss-16                                	   97382	     12475 ns/op	       0 B/op	       0 allocs/op
BenchmarkSets/runtime/map-16                                  	   56491	     21077 ns/op	      79 B/op	       0 allocs/op
```
