package a

import "testing"

// The separate test-order rule takes priority over type proximity in tests.
type testOptions struct{}

func TestFirst(t *testing.T) {}

func TestUse(t *testing.T) { _ = testOptions{} }
