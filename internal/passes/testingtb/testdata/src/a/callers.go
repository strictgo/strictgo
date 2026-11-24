package a

import "testing"

func exercise(t *testing.T) {
	helper(t)
	passInterfaces(t)
	_ = returnTB(t)
	_ = returnAny(t)
	storeInterface(t)
	ignores(t)
	_ = capture(t)
	shadowed(t)
	genericHelper[int](t, 1)
	aliasHelper(t)
	pointerAliasHelper(t)
	importAlias(t)
	TestUnregistered(t)
	inferred(t)
	needsRun(t)
	needsParallel(t)
	needsDeadline(t)
	_ = needsPointer(t)
	forwardsPointer(t)
	_ = pointerInClosure(t)
	_ = nilCheck(t)
	pointerAlias(t)
	address(t)
	dereference(t)
	assertion(t)
}
