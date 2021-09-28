# Fairy

[![Go Report Card](https://goreportcard.com/badge/github.com/rl404/fairy)](https://goreportcard.com/report/github.com/rl404/fairy)
![License: MIT](https://img.shields.io/github/license/rl404/fairy.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/rl404/fairy.svg)](https://pkg.go.dev/github.com/rl404/fairy)

_Fairy_ contains several general tools for easier and simpler development.

- Swappable cache
  - [Redis](https://redis.io/)
  - [In-memory](https://github.com/allegro/bigcache)
  - [Memcached](https://memcached.org/)
- Swappable pubsub
  - [Redis](https://redis.io/)
  - [RabbitMQ](https://rabbitmq.com/)
  - [NSQ](https://nsq.io/)
- Struct modification and validation

Hope these tools can help you or at least become a reference
for your own tools.

Good luck.

---

## Installation

```
go get github.com/rl404/fairy
```

## Quick Start

### Cache

```go
package main

import (
	"time"

	"github.com/rl404/fairy"
)

func main() {
	// Cache type.
	t := fairy.RedisCache

	// Init client.
	client, err := fairy.NewCache(t, "localhost:6379", "", time.Minute)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	key := "key"
	data := []string{"a", "b", "c"}

	// Save to cache.
	if err := client.Set(key, data); err != nil {
		panic(err)
	}

	// Get from cache.
	var newData []string
	if err := client.Get(key, &newData); err != nil {
		panic(err)
	}

	// Delete from cache.
	if err := client.Delete(key); err != nil {
		panic(err)
	}
}
```

### Pubsub

#### Publisher

```go
package main

import (
	"github.com/rl404/fairy"
)

func main() {
	// Pubsub type.
	t := fairy.RabbitMQ

	// Init client.
	client, err := fairy.NewPubSub(t, "amqp://guest:guest@localhost:5672", "")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	data := []string{"a", "b"}

	// Publish.
	if err = client.Publish("topic", data); err != nil {
		panic(err)
	}
}
```

#### Subscriber

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rl404/fairy"
)

func main() {
	// Pubsub type.
	t := fairy.RabbitMQ

	// Init client.
	client, err := fairy.NewPubSub(t, "amqp://guest:guest@localhost:5672", "")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Subscribe to specific topic/channel.
	s, err := client.Subscribe("topic")
	if err != nil {
		panic(err)
	}

	channel := s.(fairy.Channel)
	defer channel.Close()

	var newData sampleData
	msgs, errChan := channel.Read(&newData)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for {
			select {
			case <-msgs:
				// Proccess the message.
				fmt.Println(newData.Field1, newData.Field2)

			case err = <-errChan:
				// Process the error.
				fmt.Println(err)
			}
		}
	}()

	<-sigChan
}
```

### Modification and Validation

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rl404/fairy"
)

type user struct {
	Name    string `validate:"required" mod:"trim"`
	Age     int    `validate:"gt=10"`
	Magic   string `mod:"magic"`
	Country string `validate:"contain=konoha" mod:"lcase"`
}

func main() {
	// Init validator.
	v := fairy.NewValidator(true)

	// Register custom modifier.
	v.RegisterModifier("magic", func(in string) string {
		return in + " magic"
	})

	// Register custom validator.
	v.RegisterValidator("contain", func(value interface{}, param ...string) bool {
		return strings.Contains(value.(string), param[0])
	})

	// Register custom error message handler.
	v.RegisterValidatorError("gt", func(field string, param ...string) error {
		return fmt.Errorf("field %s must be greater than %s", field, param[0])
	})
	v.RegisterValidatorError("contain", func(field string, param ...string) error {
		return fmt.Errorf("field %s must contain %s", field, param[0])
	})

	// Sample 'dirty' data.
	naruto := user{
		Name:    "  Naruto ",
		Age:     15,
		Magic:   "ninja",
		Country: "Konohagakure",
	}

	// Validate struct fields.
	if err := v.Validate(&naruto); err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(naruto, "", "  ")
	fmt.Println(string(j))
}
```

*For more usage, please go to the [documentation](https://pkg.go.dev/github.com/rl404/fairy) or `example` folder.*

## License

MIT License

Copyright (c) 2021 Axel