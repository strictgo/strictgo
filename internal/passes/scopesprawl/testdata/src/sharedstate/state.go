package sharedstate // want `3 declaration groups \(count, Client, ParseConfig\); limit 2`

const defaults = 1

var count int

type Client struct{}

func (*Client) Run() { count++; println(defaults) }
func ParseConfig()   { count++; println(defaults) }
