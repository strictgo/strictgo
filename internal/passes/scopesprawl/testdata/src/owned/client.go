package owned

const defaultRetries = 2

var requests int

type Client struct{}

func NewClient() *Client     { return &Client{} }
func newClientValue() Client { return Client{} }

func (*Client) Run() {
	privateHelper()
	requests++
	println(defaultRetries)
}

func privateHelper() { innerHelper() }
func innerHelper()   {}
