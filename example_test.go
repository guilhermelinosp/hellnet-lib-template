package golanglibtemplate

import (
	"fmt"
)

// ExampleNew demonstrates the canonical New + Client.Greet happy path.
func ExampleNew() {
	c, err := New()
	if err != nil {
		panic(err)
	}
	msg, _ := c.Greet("World")
	fmt.Println(msg)
	// Output: Hello, World!
}

// ExampleNewFromEnv shows that explicit options override individual fields
// while env-first behavior is preserved: "World" is the env/default NAME and
// Repeats=2 is the explicit override.
func ExampleNewFromEnv() {
	c, err := NewFromEnv(Options{Repeats: 2})
	if err != nil {
		panic(err)
	}
	msg, _ := c.Greet("")
	fmt.Println(msg)
	// Output: Hello, World! Hello, World!
}
