# gitdoc

> Document management module using Git.

## Development

### Prerequisites

- go 1.11+

### Test

```sh
# Install dependencies.
go mod download

# Run test.
# Gitdoc directly access to file system, so it needs permission to R/W on
# `$HOME/.gitdoc` directory.
make test
```
