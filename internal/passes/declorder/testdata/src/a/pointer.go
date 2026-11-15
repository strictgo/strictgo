package a

func NewLatePointer() (*latePointer, error) { return &latePointer{}, nil }

type latePointer struct{} // want `type declarations must appear before constructors`
