package a

type value interface{ Value() int }

func identity[T value](input T) (T, error) { return input, nil }

func zero[T value]() T {
	var result T
	return result
}

func unconstrained[T any](input T) T { return input }

func integer[T ~int](input T) T { return input }

type holder[T value] struct{ input T }

func (receiver holder[T]) get() T { return receiver.input }

type valueAlias = value

func aliasedConstraint[T valueAlias](input T) T { return input }

// Explicit interface returns still hide the caller's concrete type.
func explicitInterface[T value](input T) value { return input } // want "explicitInterface returns an interface"

func mixed[T value](input T) (T, value) { return input, input } // want "mixed returns an interface"
