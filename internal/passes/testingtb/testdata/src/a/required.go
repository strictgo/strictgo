package a

import "testing"

func needsRun(t *testing.T)                { t.Run("child", func(child *testing.T) {}) }
func needsParallel(t *testing.T)           { t.Parallel() }
func needsDeadline(t *testing.T)           { _, _ = t.Deadline() }
func needsPointer(t *testing.T) *testing.T { return t }
func forwardsPointer(t *testing.T)         { needsParallel(t) }

func pointerInClosure(t *testing.T) func() *testing.T {
	return func() *testing.T { return t }
}

func nilCheck(t *testing.T) bool { return t == nil }

func pointerAlias(t *testing.T) {
	copy := t
	copy.Parallel()
}

func address(t *testing.T)     { _ = &t }
func dereference(t *testing.T) { _ = *t }
func assertion(t *testing.T)   { _ = any(t).(*testing.T) }

func namedCallback(t *testing.T)    { t.Helper() }
func assignedCallback(t *testing.T) { t.Helper() }

var (
	callbackSlot func(*testing.T) = assignedCallback
	fixedLiteral func(*testing.T) = func(t *testing.T) { t.Helper() }
)

func callbackFactory() func(*testing.T) {
	return func(t *testing.T) { t.Helper() }
}

type Different testing.T

func differentType(t *Different) { _ = t }

type (
	contract       interface{ Check(*testing.T) }
	implementation struct{}
)

// Methods may be constrained by interfaces; keep those signatures unchanged.
func (implementation) Check(t *testing.T) { t.Helper() }

var _ contract = implementation{}
