package a

import cx "context"

func importAlias(_ cx.Context) {} // want "blank parameter obscures its purpose"
