package main

type application struct{}

func main() {} // want `main entrypoints must appear before type declarations`

func NewApplication() *application { return &application{} }
