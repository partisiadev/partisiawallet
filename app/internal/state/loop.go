package state

import (
	"time"
)

func Loop() error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
		}
	}
}
