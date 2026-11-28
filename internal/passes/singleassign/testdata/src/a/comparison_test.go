package a

import "testing"

func TestComparison(t *testing.T) {
	if got, want := source(), 1; got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	got, want := source(), 1
	got, want = source(), 2 // want "use separate assignments for multiple value sources"
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if want, got := 1, source(); got != want { // want "use separate assignments for multiple value sources"
		t.Errorf("got %d, want %d", got, want)
	}
}
