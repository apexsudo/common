package sayhello

import "fmt"

// Hello returns a greeting for name, defaulting to "World" when name is empty.
func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Greet returns a custom greeting, defaulting greeting to "Hello" and name to "World".
func Greet(greeting, name string) string {
	if greeting == "" {
		greeting = "Hello"
	}
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("%s, %s!", greeting, name)
}
