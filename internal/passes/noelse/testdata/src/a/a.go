package a

func simple(condition bool) {
	if condition {
		return
	} else { // want "unnecessary else after return"
		println("work")
	}
}

func declaration(condition bool) {
	if condition {
		return
	} else { // want "unnecessary else after return"
		value := 1
		println(value)
	}
}

func empty(condition bool) {
	if condition {
		return
	} else { // want "unnecessary else after return"
	}
}

func chain(first, second bool) {
	if first {
		return
	} else if second { // want "unnecessary else after return"
		println("second")
	} else {
		println("fallback")
	}
}

func comments(condition bool) {
	if condition {
		return
	} else /* the word else in a comment */ /* preserve this comment */ { // want "unnecessary else after return"
		// Keep this body comment.
		println("work")
	} // Keep this trailing comment.
}

func nested(first, second bool) {
	if first {
		if second {
			return
		} else { // want "unnecessary else after return"
			return
		}
	} else { // want "unnecessary else after return"
		println("work")
	}
}

func nonReturning(condition bool) {
	if condition {
		println("first")
	} else {
		println("second")
	}
}

func panicking(condition bool) {
	if condition {
		panic("not handled by the initial return-only rule")
	} else {
		println("work")
	}
}
