package notanalysis

func main() {}

type namedMain struct{}

func (*namedMain) Method() {} // want `methods must appear before functions`
