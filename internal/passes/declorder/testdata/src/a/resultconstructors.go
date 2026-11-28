package a

import "bytes"

type parsedValue struct{}

func (parsedValue) validate() {}

func parseValue() (*parsedValue, error) { // want `constructors must appear before methods`
	return &parsedValue{}, nil
}

func loadValue() loadedValue { return loadedValue{} }

type loadedValue struct{} // want `type declarations must appear before constructors`

func loadBuffer() *bytes.Buffer { return &bytes.Buffer{} }

func lookupValue() (loadedValue, bool) { return loadedValue{}, false }

func useValue() loadedValue { // want `constructors must appear before functions`
	return loadedValue{}
}
