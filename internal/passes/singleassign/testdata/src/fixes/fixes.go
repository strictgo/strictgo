package fixes

func source() int { return 1 }

func split(a, b int, values []int) {
	a, b = 1, 2   // want "use separate assignments for multiple value sources"
	c, d := 3, 4  // want "use separate assignments for multiple value sources"
	a, e := 5, 6  // want "use separate assignments for multiple value sources"
	f, a := 7, 8  // want "use separate assignments for multiple value sources"
	g, _ := 9, 10 // want "use separate assignments for multiple value sources"
	_, b = 11, 12 // want "use separate assignments for multiple value sources"
	a, b = c, d   // want "use separate assignments for multiple value sources"

	// Evaluation order and dependencies prevent a safe split.
	a, b = source(), source() // want "use separate assignments for multiple value sources"
	a, b = 1, a               // want "use separate assignments for multiple value sources"
	a, b = 1, values[0]       // want "use separate assignments for multiple value sources"
	values[0], a = 1, 2       // want "use separate assignments for multiple value sources"
	a, b = 1 /* keep */, 2    // want "use separate assignments for multiple value sources"
	if a, b = 1, 2; a > b {   // want "use separate assignments for multiple value sources"
		return
	}
	for a, b = 1, 2; a < b; a, b = 2, 3 { // want "use separate assignments for multiple value sources" "use separate assignments for multiple value sources"
		break
	}
	goto label
label:
	a, b = 1, 2 // want "use separate assignments for multiple value sources"

	// A short declaration must not shadow a name used by a later source.
	{
		a, b := 1, a // want "use separate assignments for multiple value sources"
		_ = a
		_ = b
	}
	switch a {
	case 0:
		a, b = 1, 2 // want "use separate assignments for multiple value sources"
	}
	_ = b
	_ = e
	_ = f
	_ = g
}
