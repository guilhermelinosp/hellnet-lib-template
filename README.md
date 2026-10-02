# hellnet-lib-template

> GitHub template for production-ready Go libraries, pre-configured with the
> canonical Hellnet pattern: CI, linting, pre-commit hooks, dependency
> automation, and an env-first API.

Click **"Use this template"** to scaffold a new Go library in seconds.

[![pipeline](https://github.com/guilhermelinosp/hellnet-lib-template/actions/workflows/pipeline.yml/badge.svg)](https://github.com/guilhermelinosp/hellnet-lib-template/actions/workflows/pipeline.yml)
[![pr-check](https://github.com/guilhermelinosp/hellnet-lib-template/actions/workflows/pr-check.yml/badge.svg)](https://github.com/guilhermelinosp/hellnet-lib-template/actions/workflows/pr-check.yml)
[![CodeQL](https://github.com/guilhermelinosp/hellnet-lib-template/actions/workflows/codeql.yml/badge.svg)](https://github.com/guilhermelinosp/hellnet-lib-template/actions/workflows/codeql.yml)

## What's included

- **Canonical Hellnet API** — the seeded example demonstrates the exact pattern
  every Hellnet library shares (hellnet-lib-kafka, hellnet-lib-cache,
  hellnet-lib-telemetry, hellnet-lib-database, hellnet-lib-api):
  - configuration from the environment: every option is exposed as a
    `HELLNET_<LIB>_*` environment variable with a shared `HELLNET_*` fallback,
    and `.env` files load automatically (dev only, self-contained);
  - **constructors without `context.Context`** — `New`, `NewFromEnv`, `MustNew`;
    the runtime captures one `context.Background()` internally, so no public
    method takes a `ctx`;
  - options flow `DefaultOptions() → fromEnv(base) → validate()`;
  - versioning stays **v1 forever**: signature changes ship as minor/patch
    releases (never `BREAKING CHANGE`, never a major bump).
- **`.golangci.yml`** — curated linter config (errcheck, staticcheck, gosec, revive, …).
- **Lefthook** pre-commit hooks (`.lefthook.yml`): `go fmt`, `go vet`, `go mod tidy`,
  `golangci-lint`, `gitleaks`.
- **CI** (`.github/workflows`):
  - `pipeline.yml` (main): guard-semver (blocks auto-major) + semantic release +
    lib quality via [github.com/guilhermelinosp/templates].
  - `pr-check.yml` (PR): shellcheck, merge-check, gitleaks, labeler, quality.
  - `codeql.yml`: go + actions matrix.
- **Dependency automation** via Dependabot (`github-actions` + `gomod`).
- **Repo meta**: issue/PR/discussion templates, `CODEOWNERS`, `SECURITY.md`, `CONTRIBUTING.md`, `FUNDING.yml`.

## Initialize from this template

After **Use this template**, clone the new repository and run:

```bash
scripts/init-from-template.sh <repo-name> [service-name]   # renames the module, imports and cmd/ (services)
scripts/setup-repo.sh                                      # repo settings, "main" ruleset and CI variable
```

Then create the `HELLNET_ACTIONS_PRIVATE_KEY` secret (the script prints the exact command) and make sure the
`hellnet-actions` GitHub App is installed on the repository.

## Quick start

Initialise the repository first (see above); that renames the module path. Then rename the seeded API to your library:

1. Rename the package and `envPrefix` in `golanglibtemplate.go`
   (e.g. `HELLNET_KAFKA_`); keep the shared `HELLNET_` fallback.
2. Replace `Options`, `Greet` and the validation rule with your own API.
3. Keep `DefaultOptions()`, `fromEnv`, `New`/`MustNew` and `loadEnvFiles`
   — they are the canonical contract every Hellnet lib exposes.

## Configuration

| Variable | Purpose | Default |
|---|---|---|
| `HELLNET_TEMPLATE_NAME` / `HELLNET_NAME` | greeting target | `World` |
| `HELLNET_TEMPLATE_REPEATS` / `HELLNET_REPEATS` | repeat count | `1` |
| `HELLNET_TEMPLATE_VERBOSE` / `HELLNET_VERBOSE` | debug output | `false` |

`./.env` is loaded automatically by the constructors (**dev environments
only**) — no external loader call needed.

## Usage (canonical pattern)

```go
package main

import (
	"fmt"

	"github.com/<you>/<repo>"
)

func main() {
	// env-first: HELLNET_<LIB>_* (or shared HELLNET_*), .env loaded for you
	c, err := <repo>.New()
	if err != nil {
		panic(err)
	}

	// or fail fast at startup
	c = <repo>.MustNew()

	msg, err := c.Greet("World")
	if err != nil {
		panic(err)
	}
	fmt.Println(msg) // Hello, World!
}
```

## Development

```bash
go test -race ./...
go vet ./...
golangci-lint run ./...
```

Install the git hooks once with `lefthook install`: they run formatting, vet, tests (with and without `-race`), build, `go mod tidy`, lint, `govulncheck` and a secrets scan. Commits follow [Conventional Commits](https://www.conventionalcommits.org/).

## CI/CD

| Workflow | Trigger | What it does |
|---|---|---|
| `pr-check` | pull request | shellcheck, merge strategy and Conventional Commits (`merge-check`), Gitleaks, labels and the lib quality gate (module integrity, vet, race tests with coverage, lint, build, dependency review). `pr-gate` aggregates them and is the required check |
| `pipeline` | push to `main` (ignores `.github/**`) or manual | semver guard (blocks an automatic major), immutable tag + GitHub Release |
| `codeql` | nightly or manual | static analysis (CodeQL) |
| `security` | nightly or manual | Gitleaks and Trivy scans |
| `auto-pr` | push to `feat/**` or `fix/**` | opens the pull request automatically |
| `dependabot-actions-auto-merge` | Dependabot pull requests | auto-merges GitHub Actions bumps |

The workflows call reusable workflows from [templates](https://github.com/guilhermelinosp/templates) at `@latest`. Releases need the `HELLNET_ACTIONS_PRIVATE_KEY` secret and the `HELLNET_ACTIONS_CLIENT_ID` variable (set them with `scripts/setup-repo.sh`).

## Versioning

Releases and version bumps are derived from [Conventional Commits].
Hellnet libraries stay on **major v1**: a change to a public signature ships as
a minor or patch release — never a `BREAKING CHANGE`, never a major bump.

## Contributing and license

See [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md). Licensed under [Apache 2.0](LICENSE).

[github.com/guilhermelinosp/templates]: https://github.com/guilhermelinosp/templates
[Conventional Commits]: https://www.conventionalcommits.org/
