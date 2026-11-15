package blocks

const limit = 1

var defaultLimit = limit

func beforeGroups() {}

type first struct{}

func NewFirst() *first { return &first{} }

func (*first) Method() {}

func firstHelper() {}

// Starting a new type block is allowed after the first block's helper functions.
type second struct{}

func newSecond() second { return second{} }

func (second) Method() {}

func secondHelper() {}

type generic[T any] struct{ Value T }

func NewGeneric[T any]() *generic[T] { return &generic[T]{} }

func (*generic[T]) Method() {}

func genericHelper() {}
