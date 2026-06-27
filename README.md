# common

Shared packages, helpers, and utilities for ApexSudo backend services.

## Requirements

- Go 1.26+

## Installation

```bash
go get github.com/apexsudo/common
```

## Packages

| Package | Description |
|---------|-------------|
| [`pkg/config`](./pkg/config) | Generic config loader — reads JSON/TOML/YAML files into typed structs with env-var override support |
| [`pkg/ulid`](./pkg/ulid) | Prefixed, time-sortable unique identifiers (ULIDs) |
| [`pkg/workload`](./pkg/workload) | Composable server with graceful shutdown, gRPC, GraphQL, and custom workload support |
| [`pkg/graphql/dataloaders`](./pkg/graphql/dataloaders) | Generic batched dataloader constructors for solving N+1 queries in GraphQL resolvers |

## CLI

`cmd/scaffold` is a code-generation CLI for ApexSudo backend services.

### Installation

```bash
go install github.com/apexsudo/common/cmd/scaffold@latest
```

### Commands

#### `scaffold db create_table <name>`

Generates a timestamped up/down migration pair that creates a table with standard columns (`id`, `created_at`, `updated_at`, `deleted_at`) and a `deleted_at` index.

```bash
scaffold db create_table users
```

#### `scaffold db create_migration <name>`

Generates a blank timestamped up/down migration pair.

```bash
scaffold db create_migration add_email_to_users
```

Both commands walk up from the current working directory to locate `internal/storage/database/migration/scripts` and fail with a clear error if it cannot be found. The target project must have `github.com/golang-migrate/migrate/v4/cmd/migrate` registered as a Go tool (`go get -tool ...`).

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=cover.out -covermode=atomic -coverpkg=./...
```

## Releases

Releases are automated via [semantic-release](https://semantic-release.gitbook.io) on every push:

- Pushes to `main` produce a stable release (`v1.2.3`)
- Pushes to any other branch produce a pre-release tagged with the branch name (`v1.2.3-my-branch.1`)

Commit messages must follow the [Conventional Commits](https://www.conventionalcommits.org) spec to trigger a release.
