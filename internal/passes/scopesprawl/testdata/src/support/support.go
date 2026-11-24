package support // want `3 declaration groups \(Client, policy, recipe\); limit 2`

type (
	Client struct{ data record }
	record struct{}
	reader interface{ Read() }
)

type policy struct{}

func (*policy) apply() { println("independent behavior") }

type recipe struct{}

func newRecipe() *recipe { return &recipe{} }

func (*Client) Run(input reader) {
	_ = record{}
	var value policy
	value.apply()
	_ = newRecipe()
	input.Read()
}
