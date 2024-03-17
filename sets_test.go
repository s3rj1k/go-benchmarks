package main

import (
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/Workiva/go-datastructures/trie/ctrie"
	"github.com/alphadose/haxmap"
	"github.com/dghubble/trie"
	"github.com/dolthub/swiss"
	"github.com/ironpark/skiplist"
	cuckoo "github.com/panmari/cuckoofilter"

	"code.local/go-benchmarks/charhashmatrix"
	"code.local/go-benchmarks/charmatrix3d"
	"code.local/go-benchmarks/random"
)

func BenchmarkSets(b *testing.B) {
	size := math.MaxUint8

	tt := make([]string, size)
	for i := range tt {
		tt[i] = random.String(rand.Intn(size)+1, random.KubernetesNamesAllowedChars)
	}

	{
		ct := ctrie.New(nil)

		b.ResetTimer()
		b.Run("Workiva/go-datastructures/trie/ctrie", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					bytes := []byte(tt[j])

					ct.Insert(bytes, struct{}{})

					_, ok := ct.Lookup(bytes)
					if !ok {
						b.FailNow()
					}
				}

				for j := range tt {
					bytes := []byte(tt[j])

					_, _ = ct.Remove(bytes)

					_, ok := ct.Lookup(bytes)
					if ok {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		matrixHash := charhashmatrix.NewMatrix()

		b.ResetTimer()
		b.Run("local/char-xxhash-matrix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					err := matrixHash.Set(tt[j])
					if err != nil {
						b.FailNow()
					}

					if !matrixHash.Contains(tt[j]) {
						b.FailNow()
					}
				}

				for j := range tt {
					err := matrixHash.Unset(tt[j])
					if err != nil {
						b.FailNow()
					}

					if matrixHash.Contains(tt[j]) {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		matrix3D := charmatrix3d.NewMatrix(size)

		b.ResetTimer()
		b.Run("local/char-matrix-3d", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					runes := []rune(tt[j])

					err := matrix3D.Set(runes)
					if err != nil {
						b.FailNow()
					}

					if !matrix3D.Contains(runes) {
						b.FailNow()
					}
				}

				for j := range tt {
					runes := []rune(tt[j])

					err := matrix3D.Unset(runes)
					if err != nil {
						b.FailNow()
					}

					if matrix3D.Contains(runes) {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		var comp skiplist.Comparable[string] = func(lhs, rhs string) int {
			return strings.Compare(lhs, rhs)
		}

		list := skiplist.New[string, struct{}](comp)

		b.ResetTimer()
		b.Run("ironpark/skiplist", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					list.Set(tt[j], struct{}{})

					if _, ok := list.GetValue(tt[j]); !ok {
						b.FailNow()
					}
				}

				for j := range tt {
					list.Remove(tt[j])

					if _, ok := list.GetValue(tt[j]); ok {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		m := haxmap.New[string, struct{}](uintptr(size))

		b.ResetTimer()
		b.Run("alphadose/haxmap", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					m.Set(tt[j], struct{}{})

					_, ok := m.Get(tt[j])
					if !ok {
						b.FailNow()
					}
				}

				for j := range tt {
					m.Del(tt[j])

					_, ok := m.Get(tt[j])
					if ok {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		m := swiss.NewMap[string, struct{}](uint32(size))

		b.ResetTimer()
		b.Run("dolthub/swiss", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					m.Put(tt[j], struct{}{})

					if !m.Has(tt[j]) {
						b.FailNow()
					}
				}

				for j := range tt {
					m.Delete(tt[j])

					if m.Has(tt[j]) {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		cf := cuckoo.NewFilter(uint(size))

		b.ResetTimer()
		b.Run("panmari/cuckoofilter", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					bytes := []byte(tt[j])

					cf.Insert(bytes)

					ok := cf.Lookup(bytes)
					if !ok {
						b.FailNow()
					}
				}

				for j := range tt {
					bytes := []byte(tt[j])

					cf.Delete(bytes)

					ok := cf.Lookup(bytes)
					if ok {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		pt := trie.NewPathTrie()

		b.ResetTimer()
		b.Run("dghubble/trie", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					pt.Put(tt[j], struct{}{})

					val := pt.Get(tt[j])
					if val == nil {
						b.FailNow()
					}
				}

				for j := range tt {
					pt.Delete(tt[j])

					val := pt.Get(tt[j])
					if val != nil {
						b.FailNow()
					}
				}
			}
		})
	}

	{
		m := make(map[string]struct{}, size)

		b.ResetTimer()
		b.Run("runtime/map", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tt {
					m[tt[j]] = struct{}{}

					_, ok := m[tt[j]]
					if !ok {
						b.FailNow()
					}
				}

				for j := range tt {
					delete(m, tt[j])

					_, ok := m[tt[j]]
					if ok {
						b.FailNow()
					}
				}
			}
		})
	}
}
