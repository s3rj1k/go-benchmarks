package main

import (
	"math"
	"math/rand"
	"testing"

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

	matrixHash := charhashmatrix.NewMatrix()
	matrix3D := charmatrix3d.NewMatrix(size)
	set := make(map[string]struct{}, size)

	b.ResetTimer()
	b.Run("3DMatrix", func(b *testing.B) {
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

	b.ResetTimer()
	b.Run("HashMatrix", func(b *testing.B) {
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

	b.ResetTimer()
	b.Run("MapSet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := range tt {
				set[tt[j]] = struct{}{}

				_, ok := set[tt[j]]
				if !ok {
					b.FailNow()
				}
			}

			for j := range tt {
				delete(set, tt[j])

				_, ok := set[tt[j]]
				if ok {
					b.FailNow()
				}
			}
		}
	})
}
