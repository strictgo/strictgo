package a

import _ "unsafe"

const (
	first  = 1
	second = 2
)

var (
	value   = first
	another = second
)

type (
	thing          struct{}
	alias          = thing
	contract       interface{ Value() }
	number         int
	generic[T any] struct{ Value T }
)

func New() *thing                    { return &thing{} }
func NewThing() thing                { return thing{} }
func newThing() *thing               { return &thing{} }
func NewGeneric[T any]() *generic[T] { return &generic[T]{} }

func (thing) Value()     {}
func (*thing) Pointer()  {}
func (thing) NewMethod() {}

func regular() {}
func Newline() {}
func newer()   {}
func init()    {}

func localDeclarations() {
	const local = 1
	localValue := local
	type localType struct{}
	_ = localValue
	_ = localType{}
}
