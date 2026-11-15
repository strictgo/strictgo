package a

var beforeConst = 1

const ( // want `const declarations must appear before var declarations`
	lateConstant    = 1
	anotherConstant = 2
)

type config struct{}

var afterType = 1 // want `var declarations must appear before type declarations`

func NewConfig() config { return config{} }

type afterConstructor struct{}

func newAfterConstructor() afterConstructor { return afterConstructor{} }
func NewLateStruct() *lateStruct            { return &lateStruct{} }

func (config) Method() {}

var afterMethod = 1 // want `var declarations must appear before methods`

func function() {}

const afterFunction = 1       // want `const declarations must appear before functions`
var variableAfterFunction = 1 // want `var declarations must appear before functions`

type ( // want `type declarations must appear before constructors`
	lateStruct struct{}
	lateAlias  = config
)
