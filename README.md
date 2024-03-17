# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/Workiva/go-datastructures/trie/ctrie-16         	    4056	    307013 ns/op	  261921 B/op	    4904 allocs/op
BenchmarkSets/local/char-xxhash-matrix-16                     	    5190	    234906 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16                         	     992	   1224846 ns/op	  264259 B/op	     446 allocs/op
BenchmarkSets/ironpark/skiplist-16                            	    6511	    173254 ns/op	   25696 B/op	     765 allocs/op
BenchmarkSets/alphadose/haxmap-16                             	   24486	     46773 ns/op	   12240 B/op	     255 allocs/op
BenchmarkSets/dolthub/swiss-16                                	   95874	     12769 ns/op	       0 B/op	       0 allocs/op
BenchmarkSets/panmari/cuckoofilter-16                         	   24714	     45128 ns/op	   65920 B/op	     446 allocs/op
BenchmarkSets/runtime/map-16                                  	   53330	     22077 ns/op	      79 B/op	       0 allocs/op
```
