package a

type recursive struct { // want `move supporting struct "recursive" directly before parseRecursive, its first explicit use`
	next *recursive
}

func unrelatedRecursive() {}

func parseRecursive() *recursive { return &recursive{} }
