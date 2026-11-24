package generic

type Box[T any] struct{ Value T }

func NewBox[T any]() *Box[T] { return &Box[T]{} }

func (*Box[T]) Run() { helper[T]() }

func helper[T any]() {}

func Operate() { NewBox[int]().Run() }
