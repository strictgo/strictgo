package a

func counters(values []int, end int) {
	for i := end - 1; i >= 0; i-- { // want `avoid for loops with integer counters`
		println(values[i])
	}
	for i := 0; i < len(values); i++ { // want `avoid for loops with integer counters`
		println(values[i])
	}
	i := 0
	for i = 0; i < len(values); i++ { // want `avoid for loops with integer counters`
		println(values[i])
	}
	for ; i > 0; i-- { // want `avoid for loops with integer counters`
		println(i)
	}
	for i := 0; i < len(values); { // want `avoid for loops with integer counters`
		println(values[i])
		i++
	}
	for i, j := 0, end; i < j; i, j = i+1, j-1 { // want `avoid for loops with integer counters`
		println(i, j)
	}
	type count uint64
	for i := count(0); i < count(end); i++ { // want `avoid for loops with integer counters`
		println(i)
	}
	for i := int32(0); i < int32(end); i++ { // want `avoid for loops with integer counters`
		for j := 0; j < end; j++ { // want `avoid for loops with integer counters`
			println(i, j)
		}
	}
}

func allowed(values []int, ready func() bool) {
	for i, value := range values {
		println(i, value)
	}
	for i := range len(values) {
		println(i)
	}
	for ready() {
		println("waiting")
	}
	for {
		break
	}
	for value := 0.0; value < 1; value += 0.1 {
		println(value)
	}
	for text := "x"; text != ""; text = "" {
		println(text)
	}
	for _ = 0; ready(); {
		println("waiting")
	}
}
