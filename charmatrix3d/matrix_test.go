package charmatrix3d

import (
	"math/rand"
	"testing"

	"code.local/go-benchmarks/random"
)

func TestCharMatrix(t *testing.T) {
	runes := make([][]rune, 128*128)
	size := len(runes) / 4
	m := NewMatrix(size)

	for i := range runes {
		runes[i] = random.Runes(rand.Intn(size-1)+1, random.KubernetesNamesAllowedChars)

		if err := m.Set(runes[i]); err != nil {
			t.Fatalf("Set returned an error: %v", err)
		}

		if !m.Contains(runes[i]) {
			t.Fatal("Does not contain expected string after setting")
		}
	}

	for len(runes) > 0 {
		i := rand.Intn(len(runes))

		if err := m.Unset(runes[i]); err != nil {
			t.Fatalf("Unset returned an error: %v", err)
		}

		if m.Contains(runes[i]) {
			t.Fatal("Contains unexpected string after unsetting")
		}

		runes = append(runes[:i], runes[i+1:]...)
	}

	if len(runes) != 0 {
		t.Fatal("Expected that all data had been processed")
	}
}

func FuzzCharMatrix(f *testing.F) {
	for i := 0; i < 255; i++ {
		f.Add(random.String(rand.Intn(255)+1, random.KubernetesNamesAllowedChars))
	}

	f.Fuzz(func(t *testing.T, original string) {
		m := NewMatrix(255)

		runes := []rune(original)

		if err := m.Set(runes); err != nil {
			t.Fatalf("Set returned an error: %v", err)
		}

		if !m.Contains(runes) {
			t.Fatal("Does not contain expected string after setting")
		}

		if err := m.Unset(runes); err != nil {
			t.Fatalf("Unset returned an error: %v", err)
		}

		if m.Contains(runes) {
			t.Fatal("Contains unexpected string after unsetting")
		}
	})
}
