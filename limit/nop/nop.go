// Package nop is a no-operation rate limiter.
package nop

type nop struct{}

// New to create new no-operation rate limiter.
func New() *nop {
	return &nop{}
}

// Take to do nothing.
func (n *nop) Take() {}
