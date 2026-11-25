package a

// want +2 `string constant "userPath" is used only once`

const userPath = "/users"

// Names do not provide arbitrary exemptions.
// want +2 `string constant "defaultRegion" is used only once`

const defaultRegion = "us-west-2"

// want +2 `string constant "idempotencyHeader" is used only once`

const idempotencyHeader = "X-Idempotency-Key"

// want +2 `string constant "base" is used only once`

const base = "base"

// A nonliteral initializer gets a diagnostic but no overlapping fix.
// want +2 `string constant "alias" is used only once`

const alias = base

// The only reference is generated: report without offering an edit there.
// want +2 `string constant "generatedUse" is used only once`

const generatedUse = "generated consumer"

const sharedWithGenerated = "two consumers"
