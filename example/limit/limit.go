package main

import (
	"log"
	"time"

	"github.com/rl404/fairy/limit"
)

func main() {
	// Rate limit type.
	t := limit.Mutex

	// Init rate limiter.
	// For example, 5 request per second.
	l, err := limit.New(t, 5, time.Second)
	if err != nil {
		panic(err)
	}

	// Run function which is called many times.
	for i := 1; i < 13; i++ {
		// This function will be rate-limited.
		print(l, i)
	}
}

func print(l limit.Limiter, cnt int) {
	// Call Take() function for every function
	// that is rate-limited.
	l.Take()

	// Do something...
	log.Printf("run sample function: %d\n", cnt)
}
