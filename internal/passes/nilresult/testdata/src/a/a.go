package a

import "errors"

type T struct{}

func good(t *T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	return t, nil
}

func bad(t *T, err error) (*T, error) {
	return t, err // want "both returned values may be non-nil"
}

func (*T) method() (*T, error) {
	return &T{}, errors.New("failed") // want "both returned values may be non-nil"
}

func named() (t *T, err error) {
	t = &T{}
	err = errors.New("failed")
	return // want "both returned values may be non-nil"
}

func zero() (t *T, err error) { return }

func guarded(t *T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	return t, err
}

func guardedPointer(t *T, err error) (*T, error) {
	if t == nil {
		return t, err
	}
	return t, nil
}

func assigned(t *T) (*T, error) {
	var err error
	return t, err
}

func merged(t *T, b bool) (*T, error) {
	var err error
	if b {
		t = &T{}
	} else {
		t = nil
	}
	return t, err
}

func structParameter(t T, err error) (T, error) {
	return t, err // want "return a zero struct result or a nil error"
}
func ignored(n int, err error) (int, error) { return n, err }
func many(t *T, err error) (*T, int, error) { return t, 0, err }

var closure = func(t *T, err error) (*T, error) {
	return t, err // want "both returned values may be non-nil"
}

type pointerError struct{}

func (*pointerError) Error() string { return "error" }
func typedNil() (*T, error) {
	var err *pointerError
	return &T{}, err // want "both returned values may be non-nil"
}

func converted(err error) (*T, error) { return (*T)(nil), err }

func forwarded(t *T, err error) (*T, error) {
	return bad(t, err) // want "both returned values may be non-nil"
}

func nested() {
	_ = func(t *T, err error) (*T, error) {
		return t, err // want "both returned values may be non-nil"
	}
}

func namedGood() (t *T, err error) {
	t = &T{}
	return
}

type Value struct {
	N       int
	Text    string
	Flag    bool
	Pointer *T
	Nested  struct{ N int }
	Array   [2]int
	Complex complex128
}

func structGood(err error) (Value, error) {
	return Value{}, err
}

func structExplicitZero(err error) (Value, error) {
	return Value{N: 0, Text: "", Flag: false, Pointer: nil, Complex: 0}, err
}

func structBad(err error) (Value, error) {
	return Value{N: 1}, err // want "return a zero struct result or a nil error"
}

func structSuccess() (Value, error) {
	return Value{N: 1}, nil
}

func structGuarded(v Value, err error) (Value, error) {
	if err != nil {
		return Value{}, err
	}
	return v, err
}

func structNamedZero() (v Value, err error) {
	err = errors.New("failed")
	return
}

func structNamedBad() (v Value, err error) {
	v.N = 1
	err = errors.New("failed")
	return // want "return a zero struct result or a nil error"
}

func structMutated(err error) (Value, error) {
	v := Value{}
	v.Nested.N = 1
	return v, err // want "return a zero struct result or a nil error"
}

func structArray(err error) (Value, error) {
	return Value{Array: [2]int{1}}, err // want "return a zero struct result or a nil error"
}

func structForwarded(err error) (Value, error) {
	return structGood(err) // want "return a zero struct result or a nil error"
}

type OtherValue Value

func structConverted(err error) (Value, error) {
	return Value(OtherValue{}), err
}

func structMerged(b bool, err error) (Value, error) {
	var v Value
	if b {
		v = Value{}
	}
	return v, err
}

func structMergedBad(b bool, err error) (Value, error) {
	var v Value
	if b {
		v = Value{N: 1}
	}
	return v, err // want "return a zero struct result or a nil error"
}

func mutateValue(v *Value) { v.N = 1 }

func structEscaped(err error) (Value, error) {
	v := Value{}
	mutateValue(&v)
	return v, err // want "return a zero struct result or a nil error"
}

func structTypedNil() (Value, error) {
	var err *pointerError
	return Value{N: 1}, err // want "return a zero struct result or a nil error"
}

func structClosure() {
	_ = func(v Value, err error) (Value, error) {
		return v, err // want "return a zero struct result or a nil error"
	}
}

func deferredStruct(err error) (Value, error) {
	if err != nil {
		return Value{}, err
	}
	defer func() {}()
	if err != nil {
		return Value{}, err
	}
	return Value{N: 1}, nil
}

func deferredStructBad(err error) (Value, error) {
	defer func() {}()
	if err != nil {
		return Value{N: 1}, err // want "return a zero struct result or a nil error"
	}
	return Value{N: 1}, nil
}

func deferredPointer(t *T, err error) (*T, error) {
	defer func() {}()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func deferredPointerBad(t *T, err error) (*T, error) {
	defer func() {}()
	return t, err // want "both returned values may be non-nil"
}

func guardedError(t *T, err error) (*T, error) {
	if nil == err {
		return t, err
	}
	return nil, err
}
