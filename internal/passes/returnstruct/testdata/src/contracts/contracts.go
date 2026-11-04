package contracts

type Value interface{ Value() }
type Contract interface {
	Get() Value
	Close()
}
type Base interface{ Build() Value }
type Embedded interface {
	Base
	Close()
}
type Alias = Contract
