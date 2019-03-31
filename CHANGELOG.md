# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### FEATURES
- Release binaries are now automatically built and published from the CI

### ENHANCEMENTS
- Optimized Makefile

### OTHERS
- Moved CI from `Travis` to `Drone`

## [0.1.6] - 2019-03-28
### BUGFIXES
- Fixed Dockerfile build
- Fixed Travis CI builds
- Wait a bit longer for Vault container to be ready in dev-env
- Removed the IPC_LOCK capability over the build container
- Fixed the ldflags breaking darwin and windows builds

### IMPROVEMENTS
- Also build for **arm64**
- ignore dist folder in git
- Do not use go mod for build dependencies
- Tidied `go.mod`

## [0.1.5] - 2019-03-18
### FEATURES
- Added gox and ghr features to release binaries

### BUGFIXES
- Fixed a panic issue on `status` and `plan` command when the Vault path doesn't contain any value

### IMPROVEMENTS
- Updated Travis CI configuration
- Upgraded Vault to `1.0.3`
- Upgraded to golang `1.12`
- Switched to `gomodules`
- Enhanced makefile
- Updated all dependencies to their latest versions
- Made the secondary container in dev-env use the same version of Vault
- Added `IPC_LOCK` capabilities to the dev-env docker container
- Upgraded Vault libraries to `0.9.6`

### OTHERS
- Updated license to `Apache 2.0`

## [0.1.4] - 2018-02-01
### FEATURES
- Added a flag to pass sensitive content through stdin - #8
- New function `strongbox transit delete <transit_key_name>`

### BUGFIXES
- Lint CI job was failing issue since last commits
- Fixed a bug while returning an empty transit key list

### IMPROVEMENTS
- Updated dependencies
- Support Vault `0.9.3` for development env

## [0.1.3] - 2018-01-15
### FEATURES
- Embedded authentication against Vault using [approle](https://www.vaultproject.io/docs/auth/approle.html) auth backend - #6

### IMPROVEMENTS
- Switched base release container from empty (`scratch`) to `busybox` in order to be able to use it natively with GitLab CI

## [0.1.2] - 2018-01-14
### FEATURES
- Implemented a function to rotate transit keys - #2
- Got a full test environment in Makefile, added Vault container
- Possibility to generate random passwords on secret writes - #4

### IMPROVEMENTS
- Added links in Changelog
- Nicer version output
- Updated CLI, added some flags on secret write and read functions
- Enhanced functions usage outputs

### BUGFIXES
- Fixed `status` command on empty Vault Cluster

## [0.1.1] - 2018-01-13
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

[Unreleased]: https://github.com/mvisonneau/strongbox/compare/0.1.6...HEAD
[0.1.6]: https://github.com/mvisonneau/strongbox/compare/0.1.5...0.1.6
[0.1.5]: https://github.com/mvisonneau/strongbox/compare/0.1.4...0.1.5
[0.1.4]: https://github.com/mvisonneau/strongbox/compare/0.1.3...0.1.4
[0.1.3]: https://github.com/mvisonneau/strongbox/compare/0.1.2...0.1.3
[0.1.2]: https://github.com/mvisonneau/strongbox/compare/0.1.1...0.1.2
[0.1.1]: https://github.com/mvisonneau/strongbox/compare/0.1.0...0.1.1
[0.1.0]: https://github.com/mvisonneau/strongbox/tree/0.1.0
