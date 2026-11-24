package shadowing

type Client struct{}

func (*Client) Run() { helper() }
func helper()        {}

func ParseConfig() {
	helper := func() {}
	helper()
}
