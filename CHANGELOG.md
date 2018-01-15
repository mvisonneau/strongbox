# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### FEATURES
- Embedded authentication against Vault using [approle](https://www.vaultproject.io/docs/auth/approle.html) auth backend [GH-6]

### IMPROVEMENTS
- Switched base release container from empty (`scratch`) to `busybox` in order to be able to use it natively with GitLab CI

## [0.1.2] - 2018-01-15
### FEATURES
- Implemented a function to rotate transit keys [GH-2]
- Got a full test environment in Makefile, added Vault container
- Possibility to generate random passwords on secret writes [GH-4]

### IMPROVEMENTS
- Added links in Changelog
- Nicer version output
- Updated CLI, added some flags on secret write and read functions
- Enhanced functions usage outputs

### BUGFIXES
- Fixed `status` command on empty Vault Cluster

## [0.1.1] - 2018-01-14
### IMPROVEMENTS
- Added CHANGELOG.md
- Updated dependencies

### BUGFIXES
- Fixed Dockerfile
- Fixed build versioning

## [0.1.0] - 2018-01-12
### FEATURES
- Dockerfile for building the app
- Implement the CLI
- Management of state file
- Management of Vault transit keys
- Management of secrets
- Plan and apply changes on Vault
- Makefile
- CI
- Some unit tests
- License

[Unreleased]: https://github.com/mvisonneau/strongbox/compare/0.1.2...HEAD
[0.1.2]: https://github.com/mvisonneau/strongbox/compare/0.1.1...0.1.2
[0.1.1]: https://github.com/mvisonneau/strongbox/compare/0.1.0...0.1.1
[0.1.0]: https://github.com/mvisonneau/strongbox/tree/0.1.0
