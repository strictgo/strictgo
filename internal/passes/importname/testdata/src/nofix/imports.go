package nofix

import (
	foo_bar "fmt" // want `import name "foo_bar" must use lowercase`
	foobar "strings"
)

var importCollision = foo_bar.Sprint(foobar.TrimSpace(" text "))
