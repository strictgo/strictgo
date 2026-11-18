package a

type configuration[T any] struct { // want `move supporting struct "configuration" directly before parseConfiguration, its first explicit use`
	value T
}

func unrelatedGeneric() {}

func parseConfiguration[T any]() configuration[T] { return configuration[T]{} }
