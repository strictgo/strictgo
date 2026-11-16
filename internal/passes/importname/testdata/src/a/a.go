package a

import (
	json2 "encoding/json"
	Fmt "fmt"          // want `import name "Fmt" must use lowercase`
	str_util "strings" // want `import name "str_util" must use lowercase`
)

func use() {
	Fmt.Println(str_util.TrimSpace(" text "))
	_ = json2.Valid(nil)
	{
		Fmt := "a local variable, not a package reference"
		println(Fmt)
	}
}
