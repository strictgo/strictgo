package testorder

import "testing"

// A lowercase suffix is not recognized as a test by the Go testing tool.
func Testhelper(t *testing.T) {} // want "helper functions and methods must appear after test entrypoints"

func TestReal(t *testing.T) {}
