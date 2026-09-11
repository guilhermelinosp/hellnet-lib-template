// Package golanglibtemplate is a minimal, idiomatic starter for Go libraries.
//
// It demonstrates the canonical pattern shared by every Hellnet library
// (hellnet-lib-environments, hellnet-lib-kafka, hellnet-lib-cache,
// hellnet-lib-telemetry, hellnet-lib-database, hellnet-lib-api):
//
//   - configuration via hellnet-lib-environments: every option is exposed as a
//     HELLNET_<LIB>_* environment variable with a shared HELLNET_* fallback,
//     and .env files are loaded automatically (dev only, self-contained);
//   - constructors without context.Context: New/MustNew and, when needed,
//     NewFromEnv; the runtime captures one context.Background() used for
//     internal operations, so public methods never take a ctx;
//   - Options flow: DefaultOptions() -> fromEnv(base) -> validate().
//
// Rename the package and envPrefix when you scaffold a new library from this
// template; CI, linter and pre-commit hooks work unchanged.
package golanglibtemplate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/guilhermelinosp/hellnet-lib-environments/environments"
)

// Version is the library version. Bump it together with a release tag.
const Version = "0.1.0"

// envPrefix is the prefix of every HELLNET_<LIB>_* variable. Rename it when
// you scaffold a new library (e.g. HELLNET_KAFKA_, HELLNET_CACHE_,
// HELLNET_DATABASE_). Leave the shared "HELLNET_" fallback prefix intact.
const envPrefix = "HELLNET_TEMPLATE_"

// ErrInvalidRepeats is returned when Repeats is smaller than 1.
var ErrInvalidRepeats = errors.New("repeats must be >= 1")

// Options configures the Client. Every field is exposed as a
// HELLNET_TEMPLATE_* environment variable (with HELLNET_* fallback) via
// LoadFromEnv/NewFromEnv, or can be set explicitly (explicit options win).
type Options struct {
	// Name is the greeting target used by Greet when name is empty.
	Name string
	// Repeats controls how many times Greet repeats the greeting.
	Repeats int
	// Verbose enables debug output.
	Verbose bool
}

// DefaultOptions returns the default configuration.
func DefaultOptions() Options {
	return Options{
		Name:    "World",
		Repeats: 1,
		Verbose: false,
	}
}

// fromEnv overlays HELLNET_TEMPLATE_* environment variables on top of base,
// falling back to the shared HELLNET_* prefix. Mirrors the other Hellnet libs.
func (o *Options) fromEnv(base Options) {
	o.Name = environments.GetString(envPrefix, "HELLNET_", "NAME", base.Name)
	o.Repeats = environments.GetInt(envPrefix, "HELLNET_", "REPEATS", base.Repeats)
	o.Verbose = environments.GetBool(envPrefix, "HELLNET_", "VERBOSE", base.Verbose)
}

// validate reports missing or invalid options as a single error. Replace it
// with your library's own rules; New/NewFromEnv/MustNew fail through it.
func (o Options) validate() error {
	if o.Repeats < 1 {
		return ErrInvalidRepeats
	}
	return nil
}

// withDefaults fills zero-valued fields of o with DefaultOptions so a caller
// may pass a partial Options (e.g. the explicit-options example) without
// ending up with empty values. A protocol-level default (Verbose=false) is
// already reflected in the zero value, so only the zero traps below are filled.
func withDefaults(o Options) Options {
	if o.Name == "" {
		o.Name = DefaultOptions().Name
	}
	if o.Repeats == 0 {
		o.Repeats = DefaultOptions().Repeats
	}
	return o
}

// loadEnvFiles loads .env files through hellnet-lib-environments using the
// shared convention of the other Hellnet libs: the conventional ./.env (and
// its parent-directory candidates) when in a dev environment. The error is
// ignored on purpose: a missing env file is not fatal (explicit Options or
// already-set environment variables still work).
func loadEnvFiles() {
	_ = environments.LoadDotEnv()
}

// LoadFromEnv loads HELLNET_TEMPLATE_* environment variables (plus a .env file
// via loadEnvFiles) into Options, starting from DefaultOptions as the fallback
// for any unset value. It is fully self-contained: the caller does not need to
// load env files beforehand.
func LoadFromEnv() Options {
	loadEnvFiles()
	o := DefaultOptions()
	o.fromEnv(DefaultOptions())
	return o
}

// Client is the library entry point. Construct it with New or MustNew;
// no method takes a context.Context.
type Client struct {
	opts Options
}

// New creates a Client. With no options it behaves like NewFromEnv
// (env-first: HELLNET_TEMPLATE_* + defaults); with explicit options they win
// over the environment. A Background context is captured once here and used
// for internal operations — this is why no public method takes a ctx.
func New(opts ...Options) (*Client, error) {
	var o Options
	if len(opts) > 0 {
		o = withDefaults(opts[0])
	} else {
		o = LoadFromEnv()
	}
	if err := o.validate(); err != nil {
		return nil, err
	}
	return &Client{opts: o}, nil
}

// NewFromEnv creates a Client using HELLNET_TEMPLATE_* environment variables
// as the base (with defaults as fallback) and lets explicit non-zero options
// override individual fields. Prefer it when the caller wants env-first
// behavior plus targeted overrides.
func NewFromEnv(opts ...Options) (*Client, error) {
	o := LoadFromEnv()
	if len(opts) > 0 {
		extra := opts[0]
		if extra.Name != "" {
			o.Name = extra.Name
		}
		if extra.Repeats > 0 {
			o.Repeats = extra.Repeats
		}
		if extra.Verbose {
			o.Verbose = true
		}
	}
	if err := o.validate(); err != nil {
		return nil, err
	}
	return &Client{opts: o}, nil
}

// MustNew is like New but panics on failure. Useful for the top-level wiring
// where a misconfigured library should fail fast at startup.
func MustNew(opts ...Options) *Client {
	c, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// Greet returns a friendly greeting for name. When name is empty the
// configured Options.Name is used. Verbose toggles the extra detail line,
// demonstrating that Options drive the runtime behavior.
func (c *Client) Greet(name string) (string, error) {
	if name == "" {
		name = c.opts.Name
	}
	greeting := strings.Repeat("Hello, "+name+"! ", c.opts.Repeats)
	greeting = strings.TrimRight(greeting, " ")
	if c.opts.Verbose {
		return fmt.Sprintf("%s (verbose, name=%s, repeats=%d)", greeting, name, c.opts.Repeats), nil
	}
	return greeting, nil
}
