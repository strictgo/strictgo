package loopnames

func indexed(values []int) {
	for i, item := range values {
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
		println(i, item)
	}
}

func nestedRanges(rows [][]int) {
	for i, row := range rows {
		for j, item := range row {
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

func nestedCounters() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
		}
	}
}

func mixedLoops(rows [][]int) {
	for i := 0; i < len(rows); i++ {
		for j, item := range rows[i] {
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
			println(i, j, item)
		}
	}
	for i, row := range rows {
		for j := 0; j < len(row); j += 1 {
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
			println(i, j)
		}
	}
}

func indexKinds() {
	for i := range [3]int{} {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
	}
	for i := range &[3]int{} {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
	}
	for i := range "letters" {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
	}
	for i := range 3 {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
	}
}
