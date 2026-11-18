package testorder

import "testing"

const timeout = 1

var testSetting = timeout

func TestTop(t *testing.T)      {}
func BenchmarkTop(b *testing.B) {}
func FuzzTop(f *testing.F)      {}
func Example()                  {}
func TestMain(m *testing.M)     {}

type after struct{}

func NewAfter() *after { return &after{} }

func (*after) Method() {}

func helperAfter() {}
