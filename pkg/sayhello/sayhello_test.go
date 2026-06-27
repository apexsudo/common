package sayhello_test

import (
	"testing"

	"github.com/apexsudo/common/pkg/sayhello"
)

func TestHello(t *testing.T) {
	tests := []struct {
		name, input, want string
	}{
		{"with name", "Alice", "Hello, Alice!"},
		{"empty defaults to World", "", "Hello, World!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sayhello.Hello(tt.input)
			if got != tt.want {
				t.Errorf("Hello(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGreet(t *testing.T) {
	tests := []struct {
		name             string
		greeting, target string
		want             string
	}{
		{"both provided", "Hi", "Bob", "Hi, Bob!"},
		{"empty greeting defaults to Hello", "", "Bob", "Hello, Bob!"},
		{"empty name defaults to World", "Hi", "", "Hi, World!"},
		{"both empty use defaults", "", "", "Hello, World!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sayhello.Greet(tt.greeting, tt.target)
			if got != tt.want {
				t.Errorf("Greet(%q, %q) = %q, want %q", tt.greeting, tt.target, got, tt.want)
			}
		})
	}
}
