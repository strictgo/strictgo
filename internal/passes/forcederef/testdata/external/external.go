package external

func Consume(value int) {}

type Consumer struct{}

func (Consumer) Consume(value int) {}
