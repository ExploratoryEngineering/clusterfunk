package cluster

import (
	"log"
	"time"
)

// various debugging and diagnostic functions here.

func timeCall(f func(), desc string) {
	start := time.Now()
	f()
	end := time.Now()
	log.Printf("%s took %f ms", desc, float64(end.Sub(start))/float64(time.Millisecond))
}