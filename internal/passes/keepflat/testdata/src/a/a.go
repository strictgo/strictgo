package a

func boundary(b bool) {
	if b {
		for b {
			for range []int{1} {
				switch {
				case b:
					select {
					default:
						if b { // want "nesting depth 6 exceeds maximum of 5"
							if b {
							} // No repeated warning within this branch.
						}
						{ // want "nesting depth 6 exceeds maximum of 5"
						}
					}
				}
			}
		}
	}
}

func branches(b bool, x any) {
	if b {
	} else if b {
	} else if b {
	} else if b {
	} else if b {
	} else if b {
	} else {
	}
	switch x.(type) {
	case int:
		if b {
			if b {
				if b {
					if b {
						if b { // want "nesting depth 6 exceeds maximum of 5"
						}
					}
				}
			}
		}
	}
}

func closures(b bool) {
	if b {
		if b {
			if b {
				if b {
					if b {
						_ = func() {
							if b {
								if b {
									if b {
										if b {
											if b {
												if b { // want "nesting depth 6 exceeds maximum of 5"
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
