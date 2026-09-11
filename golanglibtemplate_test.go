package golanglibtemplate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	o := DefaultOptions()
	if o.Name != "World" {
		t.Fatalf("Name = %q, want %q", o.Name, "World")
	}
	if o.Repeats != 1 {
		t.Fatalf("Repeats = %d, want 1", o.Repeats)
	}
	if o.Verbose {
		t.Fatalf("Verbose = true, want false")
	}
}

func TestGreet(t *testing.T) {
	tests := []struct {
		name  string
		opts  Options
		input string
		want  string
	}{
		{name: "valid", input: "World", want: "Hello, World!"},
		{name: "uses configured name", input: "", want: "Hello, World!"},
		{name: "repeats", opts: Options{Repeats: 2}, input: "World", want: "Hello, World! Hello, World!"},
		{name: "verbose", opts: Options{Verbose: true}, input: "World", want: "Hello, World! (verbose, name=World, repeats=1)"},
		{name: "custom name via opts", opts: Options{Name: "Alice"}, input: "Bob", want: "Hello, Bob!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := withDefaults(tt.opts)
			c, err := New(o)
			if err != nil {
				t.Fatalf("New() unexpected err = %v", err)
			}
			got, err := c.Greet(tt.input)
			if err != nil {
				t.Fatalf("Greet(%q) unexpected err = %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("Greet(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewLoadsFromEnv(t *testing.T) {
	t.Setenv("HELLNET_ENVIRONMENT", "test")
	t.Setenv("HELLNET_TEMPLATE_NAME", "Env")
	t.Setenv("HELLNET_TEMPLATE_REPEATS", "2")
	t.Setenv("HELLNET_TEMPLATE_VERBOSE", "true")

	c, err := New()
	if err != nil {
		t.Fatalf("New() unexpected err = %v", err)
	}
	got, err := c.Greet("")
	if err != nil {
		t.Fatalf("Greet() unexpected err = %v", err)
	}
	want := "Hello, Env! Hello, Env! (verbose, name=Env, repeats=2)"
	if got != want {
		t.Fatalf("Greet() = %q, want %q", got, want)
	}
}

func TestNewFallsBackToSharedPrefix(t *testing.T) {
	t.Setenv("HELLNET_ENVIRONMENT", "test")
	t.Setenv("HELLNET_NAME", "Shared")
	t.Setenv("HELLNET_TEMPLATE_NAME", "")

	c, err := New()
	if err != nil {
		t.Fatalf("New() unexpected err = %v", err)
	}
	got, err := c.Greet("")
	if err != nil {
		t.Fatalf("Greet() unexpected err = %v", err)
	}
	if want := "Hello, Shared!"; got != want {
		t.Fatalf("Greet() = %q, want %q", got, want)
	}
}

func TestNewPrefersExplicitOptions(t *testing.T) {
	t.Setenv("HELLNET_ENVIRONMENT", "test")
	t.Setenv("HELLNET_TEMPLATE_NAME", "Env")

	c, err := New(Options{Name: "Explicit"})
	if err != nil {
		t.Fatalf("New() unexpected err = %v", err)
	}
	got, err := c.Greet("")
	if err != nil {
		t.Fatalf("Greet() unexpected err = %v", err)
	}
	if want := "Hello, Explicit!"; got != want {
		t.Fatalf("Greet() = %q, want %q", got, want)
	}
}

func TestNewFromEnvOverridesWithOptions(t *testing.T) {
	t.Setenv("HELLNET_ENVIRONMENT", "test")
	t.Setenv("HELLNET_TEMPLATE_NAME", "Env")
	t.Setenv("HELLNET_TEMPLATE_REPEATS", "1")

	c, err := NewFromEnv(Options{Repeats: 2})
	if err != nil {
		t.Fatalf("NewFromEnv() unexpected err = %v", err)
	}
	got, err := c.Greet("")
	if err != nil {
		t.Fatalf("Greet() unexpected err = %v", err)
	}
	if want := "Hello, Env! Hello, Env!"; got != want {
		t.Fatalf("Greet() = %q, want %q", got, want)
	}
}

func TestNewRejectsInvalidOptions(t *testing.T) {
	_, err := New(Options{Repeats: -1})
	if !errors.Is(err, ErrInvalidRepeats) {
		t.Fatalf("New() err = %v, want ErrInvalidRepeats", err)
	}
}

func TestMustNewPanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustNew() did not panic for invalid options")
		}
	}()
	MustNew(Options{Repeats: -1})
}

// TestLoadFromEnvLoadsDotEnv ensures env-first is self-contained: LoadFromEnv
// loads the conventional .env file (cwd / parents, dev only) without the caller
// having to call any external DotEnv loader. Mirrors the other Hellnet libs.
func TestLoadFromEnvLoadsDotEnv(t *testing.T) {
	t.Setenv("HELLNET_ENVIRONMENT", "test")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(
		"HELLNET_TEMPLATE_NAME=fromdotenv\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{
		"HELLNET_TEMPLATE_NAME", "HELLNET_NAME",
	} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(dir)
	defer os.Unsetenv("HELLNET_TEMPLATE_NAME")

	o := LoadFromEnv()
	if o.Name != "fromdotenv" {
		t.Fatalf("LoadFromEnv().Name = %q, want %q", o.Name, "fromdotenv")
	}
}
