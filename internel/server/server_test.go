package server

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	ticker := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-ticker.C:
			fmt.Println("111")
		}
	}
}
