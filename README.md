# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/Workiva/go-datastructures/trie/ctrie-16         	    4322	    292698 ns/op	  268371 B/op	    4924 allocs/op
BenchmarkSets/local/char-xxhash-matrix-16                     	    4944	    214875 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16                         	    1117	   1268551 ns/op	  289824 B/op	     460 allocs/op
BenchmarkSets/ironpark/skiplist-16                            	    6938	    166473 ns/op	   25693 B/op	     765 allocs/op
BenchmarkSets/alphadose/haxmap-16                             	   23404	     47650 ns/op	   12240 B/op	     255 allocs/op
BenchmarkSets/dolthub/swiss-16                                	   93891	     12623 ns/op	       0 B/op	       0 allocs/op
BenchmarkSets/panmari/cuckoofilter-16                         	   26408	     43589 ns/op	   72096 B/op	     460 allocs/op
BenchmarkSets/dghubble/trie-16                                	    4152	    278530 ns/op	  344799 B/op	    3596 allocs/op
BenchmarkSets/falmar/goradix-16                               	    3220	    388500 ns/op	  166685 B/op	    6531 allocs/op
BenchmarkSets/arriqaaq/art-16                                 	    8428	    163639 ns/op	  240304 B/op	    2301 allocs/op
BenchmarkSets/gammazero/radixtree-16                          	   14731	     82753 ns/op	   52016 B/op	    1602 allocs/op
BenchmarkSets/snorwin/gorax-16                                	    7684	    175490 ns/op	   76264 B/op	    2384 allocs/op
BenchmarkSets/armon/go-radix-16                               	   11330	    107690 ns/op	   49208 B/op	    1163 allocs/op
BenchmarkSets/runtime/map-16                                  	   45518	     24107 ns/op	      80 B/op	       0 allocs/op
```
