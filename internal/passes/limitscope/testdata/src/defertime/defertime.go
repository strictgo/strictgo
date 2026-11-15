package defertime

import "time"

func measure() {
	startedAt := time.Now() // want "pass startedAt as a deferred function argument to reduce scope"
	defer func() {
		_ = time.Since(startedAt)
	}()
}

func nested() {
	startedAt := time.Now() // want "pass startedAt as a deferred function argument to reduce scope"
	defer func() {
		report := func() { _ = time.Since(startedAt) }
		report()
	}()
}

func preservesComments() {
	// Measure how long this operation takes.
	startedAt := time.Now() // want "pass startedAt as a deferred function argument to reduce scope"
	// Report after the work completes.
	defer func() {
		_ = time.Since(startedAt)
	}()
}

func internalComment() {
	startedAt := /* keep this comment */ time.Now() // want "pass startedAt as a deferred function argument to reduce scope"
	defer func() { _ = time.Since(startedAt) }()
}

func usedOutside() {
	startedAt := time.Now()
	defer func() { _ = time.Since(startedAt) }()
	_ = startedAt
}

func reassigned() {
	startedAt := time.Now()
	defer func() { _ = time.Since(startedAt) }()
	startedAt = time.Now()
}

func notAdjacent() {
	startedAt := time.Now()
	println("intervening work")
	defer func() { _ = time.Since(startedAt) }()
}

func existingArgument() {
	startedAt := time.Now()
	defer func(value int) { _ = time.Since(startedAt); _ = value }(1)
}

func shadowedInside() {
	startedAt := time.Now()
	defer func() {
		_ = time.Since(startedAt)
		startedAt := time.Now()
		_ = time.Since(startedAt)
	}()
}

func unrelated(now func() time.Time) {
	startedAt := now()
	defer func() { _ = time.Since(startedAt) }()
}
