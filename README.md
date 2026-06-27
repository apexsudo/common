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
| [`pkg/workload`](./pkg/workload) | Composable server with graceful shutdown, gRPC, GraphQL, and custom workload support |
| [`pkg/graphql/dataloaders`](./pkg/graphql/dataloaders) | Generic batched dataloader constructors for solving N+1 queries in GraphQL resolvers |

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
