package a

func discardedAssertions(input any) {
	value, _ := input.(int)                 // want "do not discard the boolean result of a type assertion"
	value, _ = input.(int)                  // want "do not discard the boolean result of a type assertion"
	_, _ = input.(int)                      // want "do not discard the boolean result of a type assertion"
	value, _ = input.(int)                  // want "do not discard the boolean result of a type assertion"
	declared, _ := input.(int)              // want "do not discard the boolean result of a type assertion"
	if value, _ := input.(int); value > 0 { // want "do not discard the boolean result of a type assertion"
		_ = value
	}
	_ = value
	_ = declared
}

func genericAssertion[T any](input any) {
	value, _ := input.(T) // want "do not discard the boolean result of a type assertion"
	_ = value
}

func keptAssertions(input any) {
	value, ok := input.(int)
	value, ok = input.(int)
	_, ok = input.(int)
	declared, success := input.(int)
	value = input.(int)
	_ = input.(int)
	_ = value
	_ = ok
	_ = declared
	_ = success
}
