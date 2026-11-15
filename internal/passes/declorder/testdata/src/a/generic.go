package a

func NewLateGeneric[T any]() *lateGeneric[T] { return &lateGeneric[T]{} }

type lateGeneric[T any] struct{ Value T } // want `type declarations must appear before constructors`
