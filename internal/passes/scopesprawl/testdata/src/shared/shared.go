package shared // want `3 declaration groups \(Client, ParseConfig, normalize\); limit 2`

type Client struct{}

func (*Client) Run() { normalize() }

func ParseConfig() { normalize() }

func normalize() {}
