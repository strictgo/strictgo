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

type TestAlias = test.T
type PointerAlias = *test.B
type InterfaceAlias = test.TB

func aliases(t *TestAlias, b PointerAlias, f InterfaceAlias) {
	println(t)
	println(b)
	println(f)
	println("long scope")
	println("long scope")
	println("long scope")
	println("long scope")
}

func interfaces(t test.TB, b test.TB) {
	println("long scope")
	println("long scope")
	println(t)
	println(b)
	println("long scope")
	println("long scope")
	println("long scope")
}

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

type T struct{}
type Distinct test.T

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
