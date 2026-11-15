package blocks

type (
	left  struct{}
	right struct{}
)

func NewLeft() *left { return &left{} }

func (*left) Method() {}

func leftHelper() {}

// A shared type declaration establishes both groups, independently.
func NewRight() *right { return &right{} }

func (*right) Method() {}

func rightHelper() {}
