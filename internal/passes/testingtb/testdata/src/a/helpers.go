package a

import (
	"fmt"
	"testing"
)

func helper(t *testing.T) { // want `helper parameter "t" can use testing.TB`
	t.Helper()
	t.Log("message")
	t.Cleanup(func() { t.Error("cleanup") })
}

func passInterfaces(t *testing.T) { // want `helper parameter "t" can use testing.TB`
	acceptTB(t)
	acceptReporter(t)
	_ = fmt.Sprint(t)
}

func acceptTB(tb testing.TB) {}

type reporter interface{ Errorf(string, ...any) }

func acceptReporter(output reporter) {}

func returnTB(t *testing.T) testing.TB { return t } // want `helper parameter "t" can use testing.TB`

func returnAny(t *testing.T) any { return t } // want `helper parameter "t" can use testing.TB`

func storeInterface(t *testing.T) { // want `helper parameter "t" can use testing.TB`
	var common testing.TB = t
	var other testing.TB
	other = t
	_ = common
	_ = other
}

func ignores(t *testing.T) { _ = t } // want `helper parameter "t" can use testing.TB`

func capture(t *testing.T) func() { // want `helper parameter "t" can use testing.TB`
	return func() { t.Log("capture") }
}

func shadowed(t *testing.T) { // want `helper parameter "t" can use testing.TB`
	func(t *testing.T) { t.Parallel() }(nil)
	t.Helper()
}

func genericHelper[V any](t *testing.T, value V) { // want `helper parameter "t" can use testing.TB`
	t.Helper()
	_ = value
}

type (
	TestAlias    = testing.T
	PointerAlias = *testing.T
)

func aliasHelper(t *TestAlias)          { t.Helper() } // want `helper parameter "t" can use testing.TB`
func pointerAliasHelper(t PointerAlias) { t.Helper() } // want `helper parameter "t" can use testing.TB`

// This name is not a registered test because it is in a non-test source file.
func TestUnregistered(t *testing.T) { t.Helper() } // want `helper parameter "t" can use testing.TB`

var inferred = func(t *testing.T) { t.Helper() } // want `helper parameter "t" can use testing.TB`
