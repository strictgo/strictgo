package a

// Intentionally compact to check that fixes insert a statement boundary.
func oneLine(condition bool) { if condition { return } else { println("work") } } // want "unnecessary else after return"
