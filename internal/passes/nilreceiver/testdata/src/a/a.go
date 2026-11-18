package a

type receiver struct{ next *receiver }

func (r *receiver) comparisons() {
	if r == nil { // want `do not check method receiver "r" for nil`
		return
	}
	if r != nil { // want `do not check method receiver "r" for nil`
		println("present")
	}
	if nil == r { // want `do not check method receiver "r" for nil`
		return
	}
	if nil != r { // want `do not check method receiver "r" for nil`
		println("present")
	}
	if r == nil { // want `do not check method receiver "r" for nil`
		return
	}
	_ = r == nil // want `do not check method receiver "r" for nil`
	callback := func() bool {
		return r != nil // want `do not check method receiver "r" for nil`
	}
	_ = callback()
}

func (r *receiver) switches() {
	switch r {
	case nil: // want `do not check method receiver "r" for nil`
		return
	default:
		println("present")
	}
	switch {
	case r == nil: // want `do not check method receiver "r" for nil`
		return
	}
}

func (r *receiver) otherChecks(other *receiver) {
	if other == nil {
		return
	}
	if r.next == nil {
		return
	}
	if r == other {
		println("same")
	}
	switch r.next {
	case nil:
		return
	}
	alias := r
	if alias == nil {
		return
	}
	{
		r := other
		if r == nil {
			return
		}
	}
	callback := func(r *receiver) bool { return r == nil }
	_ = callback(r)
}

func (r *receiver) shadowedNil() {
	var nil *receiver
	if r == nil {
		return
	}
	switch r {
	case nil:
		return
	}
}

func (*receiver) unnamed(other *receiver) {
	if other == nil {
		return
	}
}

func ordinary(r *receiver) {
	if r == nil {
		return
	}
}

type generic[T any] struct{ Value T }

func (g *generic[T]) checked() bool {
	return g == nil // want `do not check method receiver "g" for nil`
}

type values []int

func (v values) checked() bool {
	return v == nil // want `do not check method receiver "v" for nil`
}
