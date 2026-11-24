package concepts // want `3 declaration groups \(Client, ParseConfig, RetryPolicy\); limit 2`

type Client struct{}

func NewClient() *Client { return &Client{} }
func (*Client) Run()     {}

func ParseConfig() {}

type RetryPolicy struct{}

func NewRetryPolicy() RetryPolicy { return RetryPolicy{} }
