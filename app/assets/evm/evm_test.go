package evm

import (
	"testing"
)

func BenchmarkLoadChains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := <-LoadChains()
		b.Log(ch.Done)
	}
}
