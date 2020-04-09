package main

import (
	"testing"
)

func TestGetGreeting(t *testing.T) {
	expectedGreeting := "Hello!"
	if expectedGreeting != GetGreeting() {
		t.Errorf("expected greeting to be %s, but got %s", expectedGreeting, GetGreeting())
	}
}
