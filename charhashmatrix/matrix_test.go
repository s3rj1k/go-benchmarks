package charhashmatrix

import (
	"fmt"
	"math/rand"
	"testing"

	"code.local/go-benchmarks/random"
)

func TestUint64ToHexRunes(t *testing.T) {
	const iterations = 1000

	for i := 0; i < iterations; i++ {
		val := rand.Uint64()

		result := string(uint64ToHexRunes(val))
		expected := fmt.Sprintf("%x", val)

		if result != expected {
			t.Errorf("Mismatch for value %d: expected %s, got %s", val, expected, result)
		}
	}
}

func BenchmarkToHexRunes(b *testing.B) {
	b.ResetTimer()
	b.Run("Uint64ToHex", func(b *testing.B) {
		val := rand.Uint64()

		for i := 0; i < b.N; i++ {
			_ = uint64ToHexRunes(val)
		}
	})

	b.ResetTimer()
	b.Run("FmtSprintf", func(b *testing.B) {
		val := rand.Uint64()

		for i := 0; i < b.N; i++ {
			_ = []rune(fmt.Sprintf("%x", val))
		}
	})
}

func TestHashMatrix(t *testing.T) {
	m := NewMatrix()

	tt := make([]string, 128*128)

	for i := range tt {
		length := rand.Intn(255) + 1

		tt[i] = random.String(length, random.KubernetesNamesAllowedChars)

		if err := m.Set(tt[i]); err != nil {
			t.Fatalf("Set returned an error: %v", err)
		}

		if !m.Contains(tt[i]) {
			t.Fatal("Does not contain expected string after setting")
		}
	}

	for len(tt) > 0 {
		i := rand.Intn(len(tt))

		if err := m.Unset(tt[i]); err != nil {
			t.Fatalf("Unset returned an error: %v", err)
		}

		if m.Contains(tt[i]) {
			t.Fatal("Contains unexpected string after unsetting")
		}

		tt = append(tt[:i], tt[i+1:]...)
	}

	if len(tt) != 0 {
		t.Fatal("Expected that all data had been processed")
	}
}

func FuzzHashMatrix(f *testing.F) {
	for i := 0; i < 255; i++ {
		f.Add(random.String(rand.Intn(255)+1, random.KubernetesNamesAllowedChars))
	}

	f.Fuzz(func(t *testing.T, original string) {
		m := NewMatrix()

		if err := m.Set(original); err != nil {
			t.Fatalf("Set returned an error: %v", err)
		}

		if !m.Contains(original) {
			t.Fatal("Does not contain expected string after setting")
		}

		if err := m.Unset(original); err != nil {
			t.Fatalf("Unset returned an error: %v", err)
		}

		if m.Contains(original) {
			t.Fatal("Contains unexpected string after unsetting")
		}
	})
}
