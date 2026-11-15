package main

const name = "application"

var applicationName = name

func main() { NewApplication().run() }

type application struct{}

func NewApplication() *application { return &application{} }

func (*application) run() {}

// A method named main is not the program entrypoint.
func (*application) main() {}

func helper() {}
