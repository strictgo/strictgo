package a

import bad_name "odd/path" // want `import name "bad_name" must use lowercase.*alias also repeats`

var invalidRedundant = bad_name.Value
