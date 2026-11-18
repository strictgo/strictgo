package testorder

import "testing"

func TestFirst(t *testing.T) {}

type middle struct{} // want "helper type declarations must appear after test entrypoints"

func TestSecond(t *testing.T) {}
