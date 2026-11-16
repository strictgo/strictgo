package a

type Context struct{}

func custom(_ Context) {} // want "blank parameter obscures its purpose"
