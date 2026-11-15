package blockerrors

func NewThird() third { return third{} }

func (third) Method() {}

// The type must precede its own constructor and methods, even in block mode.
type third struct{} // want `type declarations must appear before methods`
