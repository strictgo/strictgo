package tests

type Client struct{}

func (*Client) Run() { helper() }
func helper()        {}
