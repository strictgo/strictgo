package crossfile

type Client struct{}

func (*Client) Run() { privateHelper() }

func privateHelper() {}
