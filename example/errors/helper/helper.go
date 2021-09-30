package helper

import (
	"context"

	"github.com/rl404/fairy"
)

// If you are developing a project (not library),
// it's recommended to create a package or file
// like this file so you don't need to pass
// `fairy.ErrStacker` to your every functions.
// Just call `yourPkg.Init()`, `yourPkg.Wrap()`,
// or `yourPkg.Get()` and that's it.

// See `example/errors/stack.go` for example.

var stacker fairy.ErrStacker

func init() {
	stacker = fairy.NewErrStacker()
}

// Init to init context for error stack.
func Init(ctx context.Context) context.Context {
	return stacker.Init(ctx)
}

// Wrap to wrap error and put it in the stack.
func Wrap(ctx context.Context, err error, errs ...error) error {
	return stacker.Wrap(ctx, err, errs...)
}

// Get to get error stack.
func Get(ctx context.Context) interface{} {
	stacks := stacker.Get(ctx).([]string)

	// Originally, the stack order starts from the
	// deepest error but for easier debugging,
	// let's reverse it. So, it starts from
	// the outer most error. You can delete this
	// if you want the original error order.
	for i, j := 0, len(stacks)-1; i < j; i, j = i+1, j-1 {
		stacks[i], stacks[j] = stacks[j], stacks[i]
	}

	return stacks
}
