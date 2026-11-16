package nofix

func jumping(condition bool) {
	goto done
	if condition {
		return
	} else { // want "unnecessary else after return"
		local := 1
		println(local)
	}

done:
	println("done")
}
