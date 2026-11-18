package testorder

import "testing"

type service struct{} // want "helper type declarations must appear after test entrypoints"

func (service) TestMethod(t *testing.T) {} // want "helper functions and methods must appear after test entrypoints"

func TestMethodOrdering(t *testing.T) {}
