package defertime

import clock "time"

func alias() {
	startedAt := clock.Now() // want "pass startedAt as a deferred function argument to reduce scope"
	defer func() { _ = clock.Since(startedAt) }()
}

func shadowedQualifier() {
	clock := clock.Now() // want "pass clock as a deferred function argument to reduce scope"
	defer func() { println(clock.Unix()) }()
}
