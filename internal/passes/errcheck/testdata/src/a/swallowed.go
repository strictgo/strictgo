package a

func swallowed(err error) error {
	if err != nil {
		return nil // want "do not swallow a non-nil error"
	}
	return err
}

func swallowedWithValue(err error) (string, error) {
	if err != nil {
		println(err)
		return "", nil // want "do not swallow a non-nil error"
	}
	return "success", nil
}

func swallowedWithPointer(err error) (*int, error) {
	if err != nil {
		return nil, nil // want "do not swallow a non-nil error"
	}
	return nil, nil
}

func swallowedAlias(err errorAlias) errorAlias {
	if err != nil {
		return nil // want "do not swallow a non-nil error"
	}
	return err
}

func swallowedInitializer() error {
	if err := failure(); err != nil {
		return nil // want "do not swallow a non-nil error"
	}
	return nil
}

func swallowedReversed(err error) error {
	if nil != err { // want "check errors using err != nil"
		return nil // want "do not swallow a non-nil error"
	}
	return nil
}

type source struct{}

func (source) swallowed(err error) error {
	if err != nil {
		return nil // want "do not swallow a non-nil error"
	}
	return nil
}

func closures(err error) {
	callback := func() error {
		if err != nil {
			return nil // want "do not swallow a non-nil error"
		}
		return nil
	}
	_ = callback
	pointerCallback := func() *int {
		if err != nil {
			return nil
		}
		return nil
	}
	_ = pointerCallback
}

func nestedLiteral(err error) error {
	if err != nil {
		callback := func() *int { return nil }
		_ = callback
		return err
	}
	return nil
}

func swallowedMultiple(err error) (error, error) {
	if err != nil {
		return nil, nil // want "do not swallow a non-nil error"
	}
	return nil, nil
}

func pointerOnly(err error) *int {
	if err != nil {
		return nil
	}
	return nil
}

func propagated(err error) (*int, error) {
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func propagatedInAnotherPosition(err error) (error, error) {
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func reassigned(err error) error {
	if err != nil {
		err = failure()
		return nil
	}
	return err
}

func passedAddress(err error, handle func(*error)) error {
	if err != nil {
		handle(&err)
		return nil
	}
	return err
}

func shadowedError(err error) error {
	if err != nil {
		{
			err := failure()
			println(err)
		}
		return nil // want "do not swallow a non-nil error"
	}
	return err
}

func shadowedNilReturn(err error, nil error) error {
	if err != nil {
		return nil
	}
	return err
}

func bareReturn(err error) (result error) {
	if err != nil {
		return
	}
	return err
}
