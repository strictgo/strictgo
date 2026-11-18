package testorder

import "testing"

type earlyConstructed struct{} // want "helper type declarations must appear after test entrypoints"

func NewEarlyConstructed() *earlyConstructed { // want "helper functions and methods must appear after test entrypoints"
	return &earlyConstructed{}
}

func TestConstructorOrdering(t *testing.T) {}
