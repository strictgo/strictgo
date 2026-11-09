package local

func Consume(value int) {}

type Consumer struct{}

func (Consumer) Consume(value int) {}

func Generic[T any](value T) {}
