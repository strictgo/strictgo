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

func ignored(t T, err error) (T, error)     { return t, err }
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

func guardedError(t *T, err error) (*T, error) {
	if nil == err {
		return t, err
	}
	return nil, err
}
