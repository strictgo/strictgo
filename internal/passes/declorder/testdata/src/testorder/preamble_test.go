package testorder

import "testing"

func TestPreamble(t *testing.T) {}

var lateSetting = 1 // want "var declarations must appear before test entrypoints"
