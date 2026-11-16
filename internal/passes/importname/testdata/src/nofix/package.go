package nofix

import FMT "fmt" // want `import name "FMT" must use lowercase`

var (
	fmt              = 1
	packageCollision = FMT.Sprint(fmt)
)
