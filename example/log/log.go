package main

import (
	"github.com/rl404/fairy/log/chain"
)

func main() {
	// Init first logger.
	l1, err := New(Config{
		Type:       Zerolog,
		Level:      TraceLevel,
		JsonFormat: false,
		Color:      true,
	})
	if err != nil {
		panic(err)
	}

	// Init second logger.
	l2, err := New(Config{
		Type:                   Elasticsearch,
		Level:                  ErrorLevel,
		ElasticsearchAddresses: []string{"http://localhost:9200"},
		ElasticsearchUser:      "elastic",
		ElasticsearchPassword:  "",
		ElasticsearchIndex:     "logs-example",
		ElasticsearchIsSync:    true,
	})
	if err != nil {
		panic(err)
	}

	// Chain the loggers.
	l := chain.New(l1, l2)

	// General log with additional fields.
	// Key `level` can be used to differentiate
	// log level.
	l.Log(map[string]interface{}{
		"level":  ErrorLevel,
		"field1": "f1",
		"field2": "f2",
	})

	// Quick log.
	l.Trace("%s", "trace")
	l.Debug("%s", "debug")
	l.Info("%s", "info")
	l.Warn("%s", "warn")
	l.Error("%s", "error")
	// l.Fatal("%s", "fatal")
	// l.Panic("%s", "panic")

	// HTTP with log example.
	httpWithLog(l)

	// Pubsub with log example.
	pubsubWithLog(l)
}
