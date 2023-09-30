package main

import (
	"context"
	"time"
)

func main() {
	// Cache type.
	t := Redis

	// Init client.
	client, err := New(t, "localhost:6379", "", time.Minute)
	if err != nil {
		panic(err)
	}

	// Don't forget to close.
	defer client.Close()

	// Sample data. Can be any type.
	key := "key"
	data := []string{"a", "b", "c"}

	// Save to cache. Data will be encoded first.
	if err := client.Set(context.Background(), key, data); err != nil {
		panic(err)
	}

	// Create a new or use existing variable.
	// Data type should be the same as when saving to cache.
	var newData []string

	// Get data from cache. Data will be decoded to inputted
	// variable. Don't forget to use pointer.
	if err := client.Get(context.Background(), key, &newData); err != nil {
		panic(err)
	}

	// Delete data from cache.
	if err := client.Delete(context.Background(), key); err != nil {
		panic(err)
	}
}
