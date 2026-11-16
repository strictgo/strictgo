package nofix

import (
	. "strings" // want `import name "." must use lowercase`
	_ "unsafe"  // want `import name "_" must use lowercase`
)

var dotImport = TrimSpace(" text ")
