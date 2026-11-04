package a

func guard(a, b bool) int {
	switch { // want "replace boolean switch"
	case a:
		return 1
	case b:
		return 2
	default:
		return 3
	}
}

func tagged(b bool) {
	switch x := b; x { // want "replace boolean switch"
	default:
		println("default")
	case true:
		println("true")
	case false:
		println("false")
	}
}

func breaks(b bool) {
	switch b { // want "replace boolean switch"
	case true:
		break
	}
	switch { // want "replace boolean switch"
	case b:
		fallthrough
	default:
	}
	switch 1 {
	case 1:
	}
	var x any
	switch x.(type) {
	case bool:
	}
}

type flag bool

func multiple(a, b bool) {
	switch { // want "replace boolean switch"
	case a, b:
		// Keep the selected body and its comments.
		println("matched")
	}
}

func empty(b flag) {
	switch b { // want "replace boolean switch"
	}
}
