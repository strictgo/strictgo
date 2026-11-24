package a

import "testing"

func TestActual(t *testing.T) {
	t.Helper()
	helper(t)
	genericHelper[int](t, 1)
	inferred(t)
	t.Run("named", namedCallback)
	t.Run("literal", func(child *testing.T) { child.Helper() })
	callback := func(child *testing.T) { child.Helper() }
	t.Run("inferred callback", callback)
	local := func(t *testing.T) { t.Helper() } // want `helper parameter "t" can use testing.TB`
	local(t)
	func(t *testing.T) { t.Helper() }(t) // want `helper parameter "t" can use testing.TB`
}
