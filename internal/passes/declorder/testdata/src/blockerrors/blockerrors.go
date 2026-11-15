package blockerrors

type first struct{}

func NewFirst() *first { return &first{} }

func (*first) Method() {}

func firstHelper() {}

type second struct{}

func (*second) Method() {}

func NewSecond() second { return second{} } // want `constructors must appear before methods`

func secondHelper() {}

// Moving back to a previous type does not erase its ordering history.
func (*first) LateMethod() {} // want `methods must appear before functions`

func NewFirstAgain() first { return first{} } // want `constructors must appear before functions`

var lateVariable = 1 // want `var declarations must appear before functions`
