# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem . ./db/
```

## `BenchmarkSets`
```
BenchmarkSets/Workiva/go-datastructures/trie/ctrie-16         	    4269	    280659 ns/op	  265849 B/op	    4934 allocs/op
BenchmarkSets/local/char-xxhash-matrix-16                     	    6688	    167921 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-bytes-hash-matrix-16                 	    8350	    133192 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/local/char-matrix-3d-16                         	     949	   1305527 ns/op	  267045 B/op	     450 allocs/op
BenchmarkSets/ironpark/skiplist-16                            	    7489	    161048 ns/op	   25694 B/op	     765 allocs/op
BenchmarkSets/alphadose/haxmap-16                             	   26035	     46452 ns/op	   12240 B/op	     255 allocs/op
BenchmarkSets/dolthub/swiss-16                                	   88128	     14584 ns/op	       0 B/op	       0 allocs/op
BenchmarkSets/panmari/cuckoofilter-16                         	   29259	     41158 ns/op	   66912 B/op	     450 allocs/op
BenchmarkSets/dghubble/trie-16                                	    4268	    263654 ns/op	  330374 B/op	    3444 allocs/op
BenchmarkSets/falmar/goradix-16                               	    3081	    373810 ns/op	  164331 B/op	    6545 allocs/op
BenchmarkSets/arriqaaq/art-16                                 	    8426	    126893 ns/op	  213304 B/op	    2189 allocs/op
BenchmarkSets/gammazero/radixtree-16                          	   16420	     72635 ns/op	   50472 B/op	    1544 allocs/op
BenchmarkSets/snorwin/gorax-16                                	    8728	    128406 ns/op	   74184 B/op	    2318 allocs/op
BenchmarkSets/armon/go-radix-16                               	   13593	     87826 ns/op	   47744 B/op	    1105 allocs/op
BenchmarkSets/runtime/map-16                                  	   51012	     23353 ns/op	      79 B/op	       0 allocs/op
```

## `db`
```
BenchmarkSQLiteInsertSelectUpdate-16                           	   10000	    133794 ns/op	    2936 B/op	      82 allocs/op
BenchmarkSQLiteInsertSelectUpdateUsingSquirrel-16              	   10000	    150233 ns/op	   14153 B/op	     303 allocs/op
BenchmarkSQLiteInsertSelectUpdateUsingSqlf-16                  	   10000	    137016 ns/op	    4915 B/op	     148 allocs/op
BenchmarkSQLiteInsertSelectUpdateUsingTemplateWithMap-16       	   10000	    141235 ns/op	    4383 B/op	     113 allocs/op
BenchmarkSQLiteInsertSelectUpdateUsingTemplateWithStruct-16    	   10000	    137696 ns/op	    2416 B/op	      81 allocs/op
BenchmarkSQLiteInsertSelectUpdateUsingPongo2-16                	   10000	    143156 ns/op	    7366 B/op	     140 allocs/op
BenchmarkSQLiteInsertSelectUpdateUsingBuilq-16                 	   10000	    137537 ns/op	    6907 B/op	     118 allocs/op
```
