package a

func NewLateAliased() constructedAlias { return constructedAlias{} }

type lateAliased struct{} // want `type declarations must appear before constructors`

// Aliases don't define a separate main type identity.
type constructedAlias = lateAliased
