package nofix

import Type "fmt" // want `import name "Type" must use lowercase`

var keyword = Type.Sprint(1)
