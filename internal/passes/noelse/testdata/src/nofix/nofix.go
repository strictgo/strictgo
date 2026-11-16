package nofix

var value = 1

func initialized(condition bool) {
	if value := 1; condition {
		return
	} else { // want "unnecessary else after return"
		println(value)
	}
}

func redeclaration(condition bool) {
	if condition {
		return
	} else { // want "unnecessary else after return"
		local := 1
		println(local)
	}
	local := 2
	println(local)
}

func parameterCollision(condition bool, local int) {
	if condition {
		return
	} else { // want "unnecessary else after return"
		local := 1
		println(local)
	}
}

func shadowing(condition bool) {
	if condition {
		return
	} else { // want "unnecessary else after return"
		value := 2
		println(value)
	}
	println(value)
}

func closureShadowing(condition bool) {
	if condition {
		return
	} else { // want "unnecessary else after return"
		value := 2
		println(value)
	}
	callback := func() { println(value) }
	callback()
}

func innerElse(first, second bool) {
	if first {
		println("the first branch does not return")
	} else if second {
		return
	} else { // want "unnecessary else after return"
		println("must remain inside the first else")
	}
}
