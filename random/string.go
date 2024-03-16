package random

import (
	"math/rand"
)

var (
	// follow https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names with '/' as namespace/name separator
	KubernetesNamesAllowedChars = []rune("abcdefghijklmnopqrstuvwxyz0123456789-./")
)

func Runes(size int, chars []rune) []rune {
	runes := make([]rune, size)

	totalChars := len(chars)

	for i := range runes {
		runes[i] = chars[rand.Intn(totalChars)]
	}

	return runes
}

func String(size int, chars []rune) string {
	return string(Runes(size, chars))
}
