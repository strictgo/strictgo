package sharedconstants // want `3 declaration groups \(Limit, First, Second\); limit 2`

const Limit = 1

func First()  { println(Limit) }
func Second() { println(Limit) }
