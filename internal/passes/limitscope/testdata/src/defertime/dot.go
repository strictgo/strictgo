package defertime

import . "time"

func dot() {
	startedAt := Now() // want "pass startedAt as a deferred function argument to reduce scope"
	defer func() { _ = Since(startedAt) }()
}

func shadowedTime() {
	Time := 1
	startedAt := Now() // want "pass startedAt as a deferred function argument to reduce scope"
	defer func() { _ = Since(startedAt); _ = Time }()
}
