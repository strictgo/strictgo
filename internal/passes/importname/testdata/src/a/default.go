package a

import (
	"capital/path" // want `import name "Upper" must use lowercase`
	"odd/path"     // want `import name "bad_name" must use lowercase`
)

var implicit = bad_name.Value + Upper.Value
