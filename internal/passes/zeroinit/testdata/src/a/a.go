package a

import tm "time"

type Count int

const zero = 0
const disabled = false

var global = ""   // want "omit explicit zero-value initializer"
var typed int = 0 // want "omit explicit zero-value initializer"
var (
	grouped = false // want "omit explicit zero-value initializer"
)

func examples() {
	cond := false              // want "omit explicit zero-value initializer"
	text := ""                 // want "omit explicit zero-value initializer"
	n := 0                     // want "omit explicit zero-value initializer"
	f := 0.0                   // want "omit explicit zero-value initializer"
	c := 0i                    // want "omit explicit zero-value initializer"
	r := '\x00'                // want "omit explicit zero-value initializer"
	b := byte(0)               // want "omit explicit zero-value initializer"
	count := Count(0)          // want "omit explicit zero-value initializer"
	duration := tm.Duration(0) // want "omit explicit zero-value initializer"
	negative := -0.0           // want "omit explicit zero-value initializer"
	computed := 1 - 1          // want "omit explicit zero-value initializer"
	fromConst := zero          // want "omit explicit zero-value initializer"
	fromBool := disabled       // want "omit explicit zero-value initializer"
	_, _, _, _, _, _, _ = cond, text, n, f, c, r, b
	_, _, _, _, _, _ = count, duration, negative, computed, fromConst, fromBool
}

func noFix() {
	if cond := false; cond { // want "omit explicit zero-value initializer"
	}
	for i := 0; i < 1; i++ { // want "omit explicit zero-value initializer"
	}
	switch n := 0; n { // want "omit explicit zero-value initializer"
	}
	x, y := 0, ""       // want "omit explicit zero-value initializer" "omit explicit zero-value initializer"
	var a, b = false, 0 // want "omit explicit zero-value initializer" "omit explicit zero-value initializer"
	x, z := 0, 1        // x already exists: only new declarations are checked.
	_, _, _, _, _ = x, y, a, b, z
}

func ignored() {
	const flag = false
	x := 1
	s := "hello"
	b := true
	var n int
	x = 0
	var anyValue any = 0
	var ptr *int = nil
	var slice []int = nil
	var mapping map[string]int = nil
	_ = 0
	_, _, _, _, _, _, _, _ = x, s, b, n, anyValue, ptr, slice, mapping
}
