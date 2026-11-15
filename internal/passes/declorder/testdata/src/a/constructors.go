package a

// A constructor in another file does not make crossFile a main type in helpers.go.
func NewCrossFile() *crossFile { return &crossFile{} }

func NewLateValue() lateValue { return lateValue{} }

type lateValue struct{} // want `type declarations must appear before constructors`
