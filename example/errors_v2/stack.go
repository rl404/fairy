package main

import (
	"context"
	_errors "errors"
	"fmt"

	errors "github.com/rl404/fairy/errors_v2"
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
	ctx = errors.Init(ctx)

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
		return errors.Wrap(ctx, err)
	}
	return nil
}

func fn2(ctx context.Context) error {
	// Do something but return error.
	if err := fn3(ctx); err != nil {
		// Just wrap the error.
		// Add your custom error if
		// you want.
		return errors.Wrap(ctx, err, _errors.New("custom fn2 error"))
	}
	return nil
}

func fn3(ctx context.Context) error {
	// Do something but return error.
	err := _errors.New("original error")

	// But you don't want to show the error message
	// to user because, for example, it contains
	// credential. So, just wrap it and
	// add a custom error message.
	return errors.Wrap(ctx, err, _errors.New("custom fn3 error"))
}

// Create your own function to print the
// error stack.
func printStack(ctx context.Context) {
	// Get the error stacks from ctx.
	stacks := errors.Get(ctx)

	// Print however you like.
	for _, stack := range stacks {
		fmt.Println(stack.File)
		fmt.Println(stack.Function)
		fmt.Println(stack.Message)
		fmt.Println("")
	}

	// Will print from the deepest to the shallowest error:
	//
	// /fairy/example/errors_v2/stack.go:70
	// main.fn3
	// original error
	//
	// /fairy/example/errors_v2/stack.go:70
	// main.fn3
	// custom fn3 error
	//
	// /fairy/example/errors_v2/stack.go:57
	// main.fn2
	// custom fn3 error
	//
	// /fairy/example/errors_v2/stack.go:57
	// main.fn2
	// custom fn2 error
	//
	// /fairy/example/errors_v2/stack.go:46
	// main.fn1
	// custom fn2 error
	//
}
