package a

import "errors"

type errorAlias = error

func checks(err error, alias errorAlias) bool {
	if err != nil {
		println(err)
	}
	if err == nil { // want `check errors using err != nil`
		println("success")
	}
	if nil != err { // want `check errors using err != nil`
		println(err)
	}
	if nil == err { // want `check errors using err != nil`
		println("success")
	}
	if (err) == (nil) { // want `check errors using err != nil`
		println("success")
	}
	if err := failure(); err == nil { // want `check errors using err != nil`
		println("success")
	}
	if alias == nil { // want `check errors using err != nil`
		println("success")
	}
	if alias != nil {
		println(alias)
	}
	if err != nil && alias == nil { // want `check errors using err != nil`
		println(err)
	}
	return err == nil // want `check errors using err != nil`
}

func expressions(err error) {
	_ = failure() == nil // want `check errors using err != nil`
	_ = nil != failure() // want `check errors using err != nil`
	_ = failure() != nil
	var holder struct{ err error }
	_ = holder.err == nil // want `check errors using err != nil`
	_ = holder.err != nil
	_ = err == error(nil) // A conversion is not the nil identifier.
}

func failure() error { return nil }

type concreteError struct{}

func (*concreteError) Error() string { return "failure" }

func allowed(err error, other error, ptr *int, concrete *concreteError, items []int, arbitrary any) error {
	_ = ptr == nil
	_ = nil != ptr
	_ = concrete == nil
	_ = items == nil
	_ = arbitrary == nil
	_ = err == other
	_ = errors.Is(err, other)
	var target *concreteError
	_ = errors.As(err, &target)
	return err
}

func shadowedNil(err error, nil error) bool {
	return err == nil // This compares two errors, not an error with nil.
}
