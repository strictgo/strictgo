package a

import "sort"

var g int // want `variable "g" has a single-letter name outside a short scope`

const c = 1

type T struct {
	x int // Fields are not variables covered by this rule.
}

func (r *T) method() {
	println(r.x)
	println(r.x)
	println(r.x)
	println(r.x)
	println(r.x)
	println(r.x)
	println(r.x)
	println(r.x)
}

func locals() {
	x := 1    // want `variable "x" has a single-letter name outside a short scope`
	var y int // want `variable "y" has a single-letter name outside a short scope`
	π := 3    // want `variable "π" has a single-letter name outside a short scope`
	longName := 4
	_ = x
	_ = y
	_ = π
	_ = longName
	println(longName)
	println(longName)
}

func short() {
	x := 1
	println(x)
}

func boundaryAllowed() {
	x := 1
	println(x)
	println(x)
	println(x)
	println(x)
	println(x)
}

func boundaryRejected() {
	x := 1 // want `variable "x" has a single-letter name outside a short scope`
	println(x)
	println(x)
	println(x)
	println(x)
	println(x)
	println(x)
}

func loops(values []int) {
	for i := 0; i < len(values); i++ {
		println(values[i])
	}
	for i, v := range values {
		println(i, v)
	}
	for i := 0; i < len(values); i++ { // want `variable "i" has a single-letter name outside a short scope`
		println(values[i])
		println(values[i])
		println(values[i])
		println(values[i])
		println(values[i])
		println(values[i])
		println(values[i])
	}
	for i, v := range values { // want `variable "i" has a single-letter name outside a short scope` `variable "v" has a single-letter name outside a short scope`
		println(i, v)
		println(i, v)
		println(i, v)
		println(i, v)
		println(i, v)
		println(i, v)
		println(i, v)
	}
}

func parameters(x int) (y int) { // want `variable "x" has a single-letter name outside a short scope` `variable "y" has a single-letter name outside a short scope`
	y = x
	println(y)
	println(y)
	println(y)
	println(y)
	println(y)
	return
}

func shortParameters(x int) (y int) {
	return x
}

func shadowing() {
	x := 1 // want `variable "x" has a single-letter name outside a short scope`
	{
		x := 2
		println(x)
	}
	println(x)
	println(x)
	println(x)
}

func closure() {
	callback := func(x int) { // want `variable "x" has a single-letter name outside a short scope`
		println(x)
		println(x)
		println(x)
		println(x)
		println(x)
		println(x)
		println(x)
	}
	callback(1)
}

func shortBlocks() {
	if x := 1; x > 0 {
		println(x)
	}
	{
		var x int
		println(x)
	}
}

func reassignment() {
	x := 1 // want `variable "x" has a single-letter name outside a short scope`
	x = 2
	x, descriptive := 3, 4 // x is not a new definition.
	println(x)
	println(descriptive)
	println(descriptive)
	println(descriptive)
}

func comparator(values []int) {
	sort.Slice(values, func(i, j int) bool {
		a, b := values[i], values[j]
		if a != b {
			return a > b
		}
		return a < b
	})
}

var nameStr string            // want `variable "nameStr" unnecessarily contains type information with suffix "Str"`
var entriesMap map[string]int // want `variable "entriesMap" unnecessarily contains type information with suffix "Map"`
