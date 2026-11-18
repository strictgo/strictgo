package testorder

import "testing"

// The testing-specific ordering rule is restricted to _test.go files.
type productionHelper struct{}

func TestProduction(t *testing.T) {}
