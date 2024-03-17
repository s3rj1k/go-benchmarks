# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/Workiva/go-datastructures/trie/ctrie-16         	    4150	    293655 ns/op	  267730 B/op	    4924 allocs/op
BenchmarkSets/local/char-xxhash-matrix-16                     	    5114	    207552 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16                         	     966	   1122724 ns/op	  279136 B/op	     446 allocs/op
BenchmarkSets/ironpark/skiplist-16                            	    6422	    192130 ns/op	   25696 B/op	     765 allocs/op
BenchmarkSets/alphadose/haxmap-16                             	   22395	     50363 ns/op	   12240 B/op	     255 allocs/op
BenchmarkSets/dolthub/swiss-16                                	   90901	     13571 ns/op	       0 B/op	       0 allocs/op
BenchmarkSets/panmari/cuckoofilter-16                         	   24111	     46377 ns/op	   69632 B/op	     446 allocs/op
BenchmarkSets/dghubble/trie-16                                	    3843	    301771 ns/op	  332919 B/op	    3474 allocs/op
BenchmarkSets/runtime/map-16                                  	   55188	     21988 ns/op	      80 B/op	       0 allocs/op
```
