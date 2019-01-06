[![](https://img.shields.io/circleci/project/github/HelloDHLyn/gitdoc.svg?style=for-the-badge&logo=circleci&maxAge=3600)](https://circleci.com/gh/HelloDHLyn/gitdoc)
[![](https://img.shields.io/codecov/c/github/hellodhlyn/gitdoc.svg?style=for-the-badge&maxAge=3600)](https://codecov.io/gh/HelloDHLyn/gitdoc)
[![](https://img.shields.io/github/languages/top/hellodhlyn/gitdoc.svg?style=for-the-badge&colorB=375eab&maxAge=3600)](#)
[![](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge&maxAge=3600)](https://godoc.org/github.com/HelloDHLyn/gitdoc)
[![](https://img.shields.io/github/license/hellodhlyn/gitdoc.svg?style=for-the-badge&maxAge=3600)](https://opensource.org/licenses/MIT)

# gitdoc

> Library for managing documents using Git.

## Getting Started

```sh
go get github.com/hellodhlyn/gitdoc
```

See [godoc](https://godoc.org/github.com/HelloDHLyn/gitdoc) for documentation.


## Development

### Prerequisites

- go 1.11+

### Test

```sh
# Install all dependencies.
go mod download

# Run test.
# Gitdoc directly accesses to file system. For running test, it requires permissions to access `$HOME/.gitdoc` directory.
make test
```
