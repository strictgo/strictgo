package testorder

import test "testing"

func beforeTest() {} // want "helper functions and methods must appear after test entrypoints"

func TestAlias(t *test.T) {}
