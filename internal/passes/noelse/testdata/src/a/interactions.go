package a

func siblingDeclarations(first, second bool) {
	if first {
		return
	} else { // want "unnecessary else after return"
		const local = 1
	}
	if second {
		return
	} else { // want "unnecessary else after return"
		const local = 2
	}
}

func nestedScopes(first, second bool) {
	if first {
		return
	} else { // want "unnecessary else after return"
		if second {
			return
		} else { // want "unnecessary else after return"
			local := 1
			println(local)
		}
	}
	local := 2
	println(local)
}

func returningChain(first, second bool) {
	if first {
		return
	} else if second { // want "unnecessary else after return"
		return
	} else { // want "unnecessary else after return"
		println("work")
	}
}
