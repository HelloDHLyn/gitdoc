[![](https://img.shields.io/circleci/project/github/HelloDHLyn/gitdoc.svg?style=for-the-badge&logo=circleci&maxAge=3600)](https://circleci.com/gh/HelloDHLyn/gitdoc)
[![](https://img.shields.io/codecov/c/github/hellodhlyn/gitdoc.svg?style=for-the-badge&maxAge=3600)](https://codecov.io/gh/HelloDHLyn/gitdoc)
[![](https://img.shields.io/github/languages/top/hellodhlyn/gitdoc.svg?style=for-the-badge&colorB=375eab&maxAge=3600)](#)
[![](https://img.shields.io/github/license/hellodhlyn/gitdoc.svg?style=for-the-badge&maxAge=3600)](https://opensource.org/licenses/MIT)

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
