package sortnames

import sorting "sort"

func comparators(values []int) {
	sorting.Slice(values, func(i, j int) bool {
		a, b := values[i], values[j] // want `variable "a" has a single-letter name outside a short scope` `variable "b" has a single-letter name outside a short scope`
		println(a, b)
		println(a, b)
		println(a, b)
		println(a, b)
		println(a, b)
		println(a, b)
		return a < b
	})
	sorting.SliceStable(values, func(i, j int) bool {
		println(i, j)
		println(i, j)
		println(i, j)
		println(i, j)
		println(i, j)
		println(i, j)
		return values[i] < values[j]
	})
}

func nested(values []int) {
	sorting.Slice(values, func(i, j int) bool {
		{
			i := values[i] // want `variable "i" has a single-letter name outside a short scope`
			println(i)
			println(i)
			println(i)
			println(i)
			println(i)
			println(i)
		}
		callback := func(k int) { // want `variable "k" has a single-letter name outside a short scope`
			println(k)
			println(k)
			println(k)
			println(k)
			println(k)
			println(k)
			println(k)
		}
		callback(i)
		return values[i] < values[j]
	})
}

type sorter struct{}

func (sorter) Slice(_ []int, _ func(int, int) bool) {}

func unrelated(values []int, sort sorter) {
	sort.Slice(values, func(i, j int) bool { // want `variable "i" has a single-letter name outside a short scope` `variable "j" has a single-letter name outside a short scope`
		println(i, j)
		println(i, j)
		println(i, j)
		println(i, j)
		println(i, j)
		println(i, j)
		return values[i] < values[j]
	})
}
