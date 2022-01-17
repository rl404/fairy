// Package elasticsearch is a wrapper of the original "github.com/elastic/go-elasticsearch".
package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// LogLevel is level of log that will be printed.
// Will print level that is higher than your
// chosen one.
type LogLevel int8

// Available log levec.
const (
	TraceLevel LogLevel = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	Disabled
)

// Client is elasticsearch client.
type Client struct {
	client *elasticsearch.Client
	index  string
	level  LogLevel
}

// Config is elasticsearch config.
type Config struct {
	Addresses []string
	Username  string
	Password  string
	Index     string
	Level     LogLevel
}

// New to create new elasticsearch client.
func New(cfg Config) (*Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		return nil, err
	}

	// Connection check.
	if err := isError(es.Info()); err != nil {
		return nil, err
	}

	return &Client{
		client: es,
		index:  cfg.Index,
		level:  cfg.Level,
	}, nil
}

// Trace to print trace log.
func (c *Client) Trace(str string, args ...interface{}) {
	if c.level <= TraceLevel {
		c.log(TraceLevel, str, args...)
	}
}

// Debug to print debug log.
func (c *Client) Debug(str string, args ...interface{}) {
	if c.level <= DebugLevel {
		c.log(DebugLevel, str, args...)
	}
}

// Info to print info log.
func (c *Client) Info(str string, args ...interface{}) {
	if c.level <= InfoLevel {
		c.log(InfoLevel, str, args...)
	}
}

// Warn to print warn log.
func (c *Client) Warn(str string, args ...interface{}) {
	if c.level <= WarnLevel {
		c.log(WarnLevel, str, args...)
	}
}

// Error to print error log.
func (c *Client) Error(str string, args ...interface{}) {
	if c.level <= ErrorLevel {
		c.log(ErrorLevel, str, args...)
	}
}

// Fatal to print fatal log.
// Will exit the program when called.
func (c *Client) Fatal(str string, args ...interface{}) {
	if c.level <= FatalLevel {
		c.log(FatalLevel, str, args...)
		os.Exit(1)
	}
}

// Panic to print panic log.
// Will print panic error stack and exit
// like panic().
func (c *Client) Panic(str string, args ...interface{}) {
	if c.level <= PanicLevel {
		c.log(PanicLevel, str, args...)
		panic(fmt.Sprintf(str, args...))
	}
}

func (c *Client) log(level LogLevel, str string, args ...interface{}) {
	c.Log(map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"level":   level,
		"message": fmt.Sprintf(str, args...),
	})
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (c *Client) Log(fields map[string]interface{}) {
	if len(fields) == 0 {
		return
	}

	if _, ok := fields["time"]; !ok {
		fields["time"] = time.Now().Format(time.RFC3339)
	}

	if level, ok := fields["level"]; ok {
		lvl := LogLevel(reflect.ValueOf(level).Int())
		fields["level"] = c.getLevelStr(lvl)
		if c.level <= lvl {
			c.send(fields)
		}
		return
	}

	c.send(fields)
}

func (c *Client) send(data interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	req := esapi.IndexRequest{
		Index:   c.index,
		Body:    strings.NewReader(string(d)),
		Refresh: "true",
	}
	if err := isError(req.Do(context.Background(), c.client)); err != nil {
		log.Println(err)
	}
}

func isError(res *esapi.Response, err error) error {
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	return nil
}

func (c *Client) getLevelStr(lvl LogLevel) string {
	switch lvl {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	default:
		return "log"
	}
}