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
| [`pkg/sayhello`](./pkg/sayhello) | Greeting utilities |

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
