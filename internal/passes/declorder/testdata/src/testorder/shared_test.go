package testorder

// A file containing only shared helpers does not require a test function.
type shared struct{}

func (shared) Method() {}

func sharedFunction() {}
