package analysis

type Pass struct{}

type Analyzer struct {
	Name string
	Run  func(*Pass) (any, error)
}
