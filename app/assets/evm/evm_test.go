package evm

import (
	"testing"
)

func BenchmarkLoadChainsFromAssets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LoadChains()
	}
}
