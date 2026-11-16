package recovernames

func conditional() {
	if r := recover(); r != nil {
		println(r)
		println(r)
		println(r)
		println(r)
		println(r)
		println(r)
		println(r)
	}
}

func assigned() {
	r := recover()
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
}

func declared() {
	r := recover()
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
}

func shadowedBuiltin(recover func() any) {
	if r := recover(); r != nil { // want `variable "r" has a single-letter name outside a short scope`
		println(r)
		println(r)
		println(r)
		println(r)
		println(r)
		println(r)
		println(r)
	}
}

func reassigned() {
	var r any // want `variable "r" has a single-letter name outside a short scope`
	r = recover()
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
	println(r)
}

func shadowedVariable() {
	if r := recover(); r != nil {
		{
			r := 1 // want `variable "r" has a single-letter name outside a short scope`
			println(r)
			println(r)
			println(r)
			println(r)
			println(r)
			println(r)
			println(r)
		}
		println(r)
	}
}
