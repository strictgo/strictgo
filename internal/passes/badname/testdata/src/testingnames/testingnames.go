package testingnames

import test "testing"

func conventional(t *test.T, b *test.B, f *test.F, m *test.M) {
	println(t)
	println(b)
	println(f)
	println(m)
	println("long scope")
	println("long scope")
	println("long scope")
}

type (
	TestAlias      = test.T
	PointerAlias   = *test.B
	InterfaceAlias = test.TB
)

func aliases(t *TestAlias, b PointerAlias, f InterfaceAlias) {
	println(t)
	println(b)
	println(f)
	println("long scope")
	println("long scope")
	println("long scope")
	println("long scope")
}

func interfaces(t test.TB, b test.TB) { // want `testing.TB variable "t" should be named tb`
	println("long scope")
	println("long scope")
	println(t)
	println(b)
	println("long scope")
	println("long scope")
	println("long scope")
}

func shortInterface(t test.TB) { t.Helper() } // want `testing.TB variable "t" should be named tb`

func aliasedInterface(t InterfaceAlias) {} // want `testing.TB variable "t" should be named tb`

type testState struct {
	t test.TB // want `testing.TB variable "t" should be named tb`
}

func commonInterface(tb test.TB) { tb.Helper() }

type TB interface{ Helper() }

func unrelatedInterface(t TB) { t.Helper() }

func subtest(parent *test.T) {
	parent.Run("child", func(t *test.T) {
		t.Log("one")
		t.Log("two")
		t.Log("three")
		t.Log("four")
		t.Log("five")
		t.Log("six")
		t.Log("seven")
	})
}

type (
	T        struct{}
	Distinct test.T
)

func unrelated(t *T, b *Distinct) { // want `variable "t" has a single-letter name outside a short scope` `variable "b" has a single-letter name outside a short scope`
	println(t)
	println(b)
	println("long scope")
	println("long scope")
	println("long scope")
	println("long scope")
	println("long scope")
}

func nonconventional(x *test.T, t test.T, b *test.T) { // want `variable "x" has a single-letter name outside a short scope` `variable "t" has a single-letter name outside a short scope` `variable "b" has a single-letter name outside a short scope`
	println(x)
	println(&t)
	println(b)
	println("long scope")
	println("long scope")
	println("long scope")
	println("long scope")
}
