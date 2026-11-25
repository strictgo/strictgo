package a

const (
	StatusPending        = "pending"
	explicit      string = "typed"
)

type status string

const pending status = "pending"

const (
	grouped = "grouped"
)

const (
	first, second = "first", "second"
	unused        = "unused"
	shared        = "shared"
	number        = 42
)

const query = `
SELECT id, name
FROM users
WHERE id = $1
`

const (
	escaped = "first\nsecond"
	joined  = "first" +
		"second"
)

func preserved() []string {
	return []string{StatusPending, explicit, string(pending), grouped, first, second, query, escaped, joined, shared, shared}
}

func nonString() int { return number }
