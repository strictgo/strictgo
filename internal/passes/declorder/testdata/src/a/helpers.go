package a

// A support type does not force subsequent variables into the wrong section.
type supportConfig struct{}

var supportValue = 1

func helperFunction() {}

// Supporting declarations may stay near their uses, including the visit case.
type visit struct{ value int }

func useVisit() int { return visit{}.value }

type crossFile struct{}

type supportAlias = supportConfig

type supportGeneric[T any] struct{ Value T }
