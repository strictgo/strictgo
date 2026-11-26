package nofix

import (
	. "strings" // want `import name "." must use lowercase`
	_ "unsafe"
)

var dotImport = TrimSpace(" text ")
