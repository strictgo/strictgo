package functions // want `3 declaration groups \(first, second, normalize\); limit 2`

func first()     { normalize() }
func second()    { normalize() }
func normalize() {}
