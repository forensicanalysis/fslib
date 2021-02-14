// Package context provides an simple replacement for the standard library
// context for older Go versions.
package context

type Context interface {
	Done() <-chan struct{}
}

type emptyCtx int

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func Background() Context {
	return new(emptyCtx)
}
