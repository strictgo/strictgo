package a

import (
	"context"
	"testing"
)

type ContextAlias = context.Context

type HTTPClient struct{}

func ordinary(_ context.Context) {} // want "blank parameter obscures its purpose"

func aliased(_ ContextAlias) {} // want "blank parameter obscures its purpose"

func (HTTPClient) Method(_ context.Context) {} // want "blank parameter obscures its purpose"

var literal = func(_ context.Context) {} // want "blank parameter obscures its purpose"

type Handler func(_ context.Context) // want "blank parameter obscures its purpose"

type HandlerAlias = func(_ context.Context) // want "blank parameter obscures its purpose"

type Service interface {
	Handle(_ context.Context) // want "blank parameter obscures its purpose"
}

type container struct {
	callback func(_ context.Context) // want "blank parameter obscures its purpose"
}

var callable func(_ context.Context) // want "blank parameter obscures its purpose"

var converted = (func(_ context.Context))(nil) // want "blank parameter obscures its purpose"

func accepts(callback func(_ context.Context)) {} // want "blank parameter obscures its purpose"

func returns() func(_ context.Context) { return nil } // want "blank parameter obscures its purpose"

func local() {
	var callback func(_ context.Context) // want "blank parameter obscures its purpose"
	_ = callback
}

func generic[T any](_ context.Context, _ T) {} // want "blank parameter obscures its purpose" "blank parameter obscures its purpose"

func client(_ *HTTPClient) {} // want "blank parameter obscures its purpose"

func basics(_ string, _ bool, _ int, _ error) {} // want "blank parameter obscures its purpose" "blank parameter obscures its purpose" "blank parameter obscures its purpose" "blank parameter obscures its purpose"

func variadic(_ ...string) {} // want "blank parameter obscures its purpose"

func testCallback(_ *testing.T) {} // want "blank parameter obscures its purpose"

// Unnamed parameters and blank results are outside this rule.
func unnamed(context.Context) {}
func blankResult() (_ int)    { return 0 }

// Field names do not collide with a newly introduced parameter.
func fieldUse(_ context.Context) { // want "blank parameter obscures its purpose"
	var record struct{ ctx int }
	_ = record.ctx
}
