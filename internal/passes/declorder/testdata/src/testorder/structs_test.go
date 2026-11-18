package testorder

import "testing"

type fixture struct{} // want "helper type declarations must appear after test entrypoints"

func TestAfterStruct(t *testing.T) {}
