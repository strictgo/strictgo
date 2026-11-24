package sharedconstants

const defaults = 1

type Client struct{}

func (*Client) Run() { println(defaults) }
func ParseConfig()   { println(defaults) }
