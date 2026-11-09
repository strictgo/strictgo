package fixes

import "strings"

type (
	run     struct{ Value int }
	rebuild struct{ Run run }
	capture struct{ Rebuild *rebuild }
	record  struct {
		Status string
		Stage  string
		Number int
	}
)

func expanded(normal, noVCS *capture, values []int, table map[int]int, result *record) {
	base, variant := normal.Rebuild.Run, noVCS.Rebuild.Run                // want "use separate assignments for multiple value sources"
	left, right := values[0], values[1]                                   // want "use separate assignments for multiple value sources"
	before, after := table[0], table[1]                                   // want "use separate assignments for multiple value sources"
	first, second := source(), source()                                   // want "use separate assignments for multiple value sources"
	seen, recompiled := make(map[string]bool), make(map[string]bool)      // want "use separate assignments for multiple value sources"
	low, high := -1, 1+2                                                  // want "use separate assignments for multiple value sources"
	text, repeated := strings.TrimSpace(" text "), strings.Repeat("x", 2) // want "use separate assignments for multiple value sources"
	one, two := run{Value: 1}, run{Value: 2}                              // want "use separate assignments for multiple value sources"
	function, number := func() int { return 1 }, 2                        // want "use separate assignments for multiple value sources"
	Run, alternate := normal.Rebuild.Run, noVCS.Rebuild.Run               // want "use separate assignments for multiple value sources"
	result.Status, result.Stage = "failed", "cleanup"                     // want "use separate assignments for multiple value sources"
	result.Number, high = first, -2                                       // want "use separate assignments for multiple value sources"
	left, right = one.Value, two.Value                                    // want "use separate assignments for multiple value sources"
	left, right = -3, 2+3                                                 // want "use separate assignments for multiple value sources"

	// Existing variables can be observed by calls or by panic recovery.
	left, extra := 1, source()                // want "use separate assignments for multiple value sources"
	left, right = 1, values[0]                // want "use separate assignments for multiple value sources"
	left, right = 1, normal.Rebuild.Run.Value // want "use separate assignments for multiple value sources"
	// The pointer may alias the struct read by the second source.
	var copied record
	var other record
	result.Number, other = 1, copied // want "use separate assignments for multiple value sources"
	_ = other
	// Splitting must not bind a later expression to a new local name.
	{
		strings, other := 1, strings.TrimSpace(" text ") // want "use separate assignments for multiple value sources"
		_ = strings
		_ = other
	}
	_ = base
	_ = variant
	_ = left
	_ = right
	_ = before
	_ = after
	_ = first
	_ = second
	_ = seen
	_ = recompiled
	_ = low
	_ = high
	_ = text
	_ = repeated
	_ = one
	_ = two
	_ = function
	_ = number
	_ = Run
	_ = alternate
	_ = extra
}
