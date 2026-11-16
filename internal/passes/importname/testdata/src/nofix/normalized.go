package nofix

import (
	Foo "fmt"      // want `import name "Foo" must use lowercase`
	_foo "strings" // want `import name "_foo" must use lowercase`
)

var normalizedCollision = Foo.Sprint(_foo.TrimSpace(" text "))
