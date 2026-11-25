package a

import "fmt"

func createUser() error {
	// want +2 `string constant "errorMessage" is used only once`

	const errorMessage = "failed to create user"
	return fmt.Errorf(errorMessage)
}

func escapedQuote() string {
	// want +2 `string constant "quoted" is used only once`

	const quoted = `say "hello"\world`
	return quoted
}

const shadow = "unused outer"

func shadowed() string {
	// want +2 `string constant "shadow" is used only once`

	const shadow = "inner"
	return shadow
}

func parenthesized() string {
	// want +2 `string constant "message" is used only once`

	const message = "parenthesized"
	return message
}
