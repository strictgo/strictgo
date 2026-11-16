package nofix

import Any "fmt" // want `import name "Any" must use lowercase`

var builtin = Any.Sprint(1)
