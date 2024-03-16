# go-benchmarks

## How to run benchmarks
```
go test -bench=. -benchmem
```

## `sets_test.go`
```
BenchmarkSets/3DMatrix-16         	     970	   1141171 ns/op	  275072 B/op	     464 allocs/op
BenchmarkSets/HashMatrix-16       	    7248	    160612 ns/op	   65280 B/op	    1020 allocs/op
BenchmarkSets/MapSet-16           	   56374	     21092 ns/op	      74 B/op	       0 allocs/op
```
