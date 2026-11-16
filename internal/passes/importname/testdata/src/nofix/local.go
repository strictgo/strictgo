package nofix

import TIME "time" // want `import name "TIME" must use lowercase`

func localCollision() {
	time := 1
	println(time, TIME.Now().Unix())
}
