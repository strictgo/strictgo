package a

type methodBeforeConstructor struct{}

func (methodBeforeConstructor) Act() {}

func NewLate() *methodBeforeConstructor { // want `constructors must appear before methods`
	return &methodBeforeConstructor{}
}

func ordinary() {}

func (*methodBeforeConstructor) Later() {} // want `methods must appear before functions`

func newLate() methodBeforeConstructor { // want `constructors must appear before functions`
	return methodBeforeConstructor{}
}
