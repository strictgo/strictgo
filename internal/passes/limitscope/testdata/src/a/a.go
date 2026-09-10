package a

import (
	"errors"
	"fmt"
)

func base() error {
	err := f() // want `inline err assignment and check to reduce scope`
	if err != nil {
		return err
	}
	return nil
}

func inlinedAndNotInlined() error {
	if err := f(); err != nil {
		return err
	}
	err := f() // want `inline err assignment and check to reduce scope`
	if err != nil {
		return err
	}
	return nil
}

func multipleResultParameter() error {
	s, err := g()
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func f() error {
	return errors.New("failed")
}

func g() (string, error) {
	return "", errors.New("failed")
}
