package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rl404/fairy/example/errors/helper"
)

func main() {
	// Call some function.
	fn()
}

func fn() {
	// Create a context.
	// You have to pass this to every function
	// called in fn() so you can list all
	// the wrapped errors.
	ctx := context.Background()

	// Init the stack.
	// Put this once in the outer most of
	// your function. Like this function.
	ctx = helper.Init(ctx)

	// Optional.
	// Print the stacked errors.
	defer printStack(ctx)

	// Call some function with initiated context
	// in the param.
	if err := fn1(ctx); err != nil {
		fmt.Println(err)
	}
}

func fn1(ctx context.Context) error {
	// You don't need to init the stack like
	// in fn() because it's already initiated.

	// Do something but return error.
	if err := fn2(ctx); err != nil {
		// Just wrap the error.
		return helper.Wrap(ctx, err)
	}
	return nil
}

func fn2(ctx context.Context) error {
	// Do something but return error.
	if err := fn3(ctx); err != nil {
		// Just wrap the error.
		// Add your custom error if
		// you want.
		return helper.Wrap(ctx, errors.New("custom fn2 error"), err)
	}
	return nil
}

func fn3(ctx context.Context) error {
	// Do something but return error.
	err := errors.New("original error")

	// But you don't want to show the error message
	// to user because, for example, it contains
	// credential. So, just wrap it and
	// add a custom error message.
	return helper.Wrap(ctx, errors.New("custom error fn3"), err)
}

// Create your own function to print the
// error stack.
func printStack(ctx context.Context) {
	// Convert the stack to whatever your tool
	// implement the interface. In this case,
	// []string.
	stacks := helper.Get(ctx).([]string)

	// Format however you like and print it.
	if len(stacks) > 0 {
		fmt.Println(strings.Join(stacks, "\n"))
	}

	// Will print:
	// stack.go:47
	// stack.go:58 custom fn2 error
	// stack.go:58
	// stack.go:71 custom error fn3
	// stack.go:71 original error
}
