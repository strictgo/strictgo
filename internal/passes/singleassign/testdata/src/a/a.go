package a

func pair() (int, int) { return 1, 2 }
func source() int      { return 1 }

func assignments(a, b, c int, values []int, channel <-chan int) {
	a, b = 1, 2                               // want "use separate assignments for multiple value sources"
	a, b, c = 1, 2, 3                         // want "use separate assignments for multiple value sources"
	x, y := 1, 2                              // want "use separate assignments for multiple value sources"
	a, b = source(), source()                 // want "use separate assignments for multiple value sources"
	a, b = b, 1                               // want "use separate assignments for multiple value sources"
	a, b = a, b                               // want "use separate assignments for multiple value sources"
	a, b = b, b                               // want "use separate assignments for multiple value sources"
	a, a = b, a                               // want "use separate assignments for multiple value sources"
	a, _ = b, a                               // want "use separate assignments for multiple value sources"
	x, y = y, 1                               // want "use separate assignments for multiple value sources"
	values[0], values[1] = 1, 2               // want "use separate assignments for multiple value sources"
	values[source()], a = a, values[source()] // want "use separate assignments for multiple value sources"
	a, b = <-channel, <-channel               // want "use separate assignments for multiple value sources"
	a, b = 1,                                 // want "use separate assignments for multiple value sources"
		2

	a, b = b, a
	a, b, c = b, c, a
	a, b, c = b, a, c
	a, b = b, a
	values[0], values[1] = values[1], values[0]
	left := struct{ a, b int }{}
	left.a, left.b = left.b, left.a
	a, b = pair()
	x, ok := map[int]int{}[0]
	_, present := any(x).(int)
	x, open := <-channel
	x = 1
	_ = ok
	_ = present
	_ = open
	_ = y
	_ = c
}
