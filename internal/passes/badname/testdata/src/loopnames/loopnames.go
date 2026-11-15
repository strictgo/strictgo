package loopnames

var i int // want `variable "i" has a single-letter name outside a short scope`

func whileStyle() {
	var i int
	for {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		i++
		if i > 10 {
			break
		}
	}
}

func conditional() {
	var i int
	for i < 10 {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		i += 2
	}
}

func rangeBody(values []int) {
	var i int
	for range values {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		i++
	}
}

func shadowed() {
	var i int // want `variable "i" has a single-letter name outside a short scope`
	{
		var i int
		for {
			i++
			println(i)
			println(i)
			println(i)
			println(i)
			println(i)
			println(i)
			break
		}
	}
	println(i)
}

func callbackOnly() {
	var i int // want `variable "i" has a single-letter name outside a short scope`
	for {
		callback := func() { i++ }
		callback()
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		break
	}
}

func notInLoop() {
	var i int // want `variable "i" has a single-letter name outside a short scope`
	println(i)
	println(i)
	println(i)
	println(i)
	println(i)
	println(i)
	i++
}

func nonInteger() {
	var i float64 // want `variable "i" has a single-letter name outside a short scope`
	for {
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		println(i)
		i++
		break
	}
}

func otherName() {
	var j int // want `variable "j" has a single-letter name outside a short scope`
	for {
		println(j)
		println(j)
		println(j)
		println(j)
		println(j)
		println(j)
		j++
		break
	}
}

func globalCounter() {
	for {
		i++
		break
	}
}
