package a

import test "testing"

func importAlias(t *test.T) { t.Helper() } // want `helper parameter "t" can use testing.TB`
