package nofix

// A deliberate explanation must not be discarded or reattached.
const documented = "documented" // want `string constant "documented" is used only once`

const trailing = "trailing" /* important explanation */ // want `string constant "trailing" is used only once`

const embedded = /* explanation */ "embedded" // want `string constant "embedded" is used only once`

// want +2 `string constant "computed" is used only once`

const computed = "left" + "right"

func use() []string { return []string{documented, trailing, embedded, computed} }
