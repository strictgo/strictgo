package loopnames

func standaloneJ(values []int) {
	for j, item := range values { // want `variable "j" has a single-letter name outside a short scope`
		println(j, item)
		println(j, item)
		println(j, item)
		println(j, item)
		println(j, item)
		println(j, item)
		println(j, item)
	}
}

func mapKey(values map[int]string) {
	for i, item := range values { // want `variable "i" has a single-letter name outside a short scope`
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
	}
}

func channelValue(values <-chan int) {
	for i := range values { // want `variable "i" has a single-letter name outside a short scope`
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
	}
}

func iteratorValue(values func(func(int, string) bool)) {
	for i, item := range values { // want `variable "i" has a single-letter name outside a short scope`
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
	}
}

func sliceValue(values []int) {
	for _, i := range values { // want `variable "i" has a single-letter name outside a short scope`
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
	}
}

func reassignedIndex(values []int) {
	var i int // want `variable "i" has a single-letter name outside a short scope`
	for i = range values {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
	}
}

func callbackIndex(rows [][]int) {
	for i, row := range rows {
		callback := func() {
			for j, item := range row { // want `variable "j" has a single-letter name outside a short scope`
				println(i, j, item)
				println(i, j, item)
				println(i, j, item)
				println(i, j, item)
				println(i, j, item)
				println(i, j, item)
				println(i, j, item)
			}
		}
		callback()
	}
}

func wrongOuterIndex(rows [][]int) {
	for rowIndex, row := range rows {
		for j, item := range row { // want `variable "j" has a single-letter name outside a short scope`
			println(rowIndex, j, item)
			println(rowIndex, j, item)
			println(rowIndex, j, item)
			println(rowIndex, j, item)
			println(rowIndex, j, item)
			println(rowIndex, j, item)
			println(rowIndex, j, item)
		}
	}
}

func mapOuterIndex(rows map[int][]int) {
	for i, row := range rows { // want `variable "i" has a single-letter name outside a short scope`
		for j, item := range row { // want `variable "j" has a single-letter name outside a short scope`
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
		}
	}
}
