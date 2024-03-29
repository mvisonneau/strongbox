# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to  [0ver](https://0ver.org).

## [Unreleased]

## [v0.2.2] - 2022-02-11

### Added

- Release container images to `quay.io`

### Changed

- Bumped go to `1.17`
- Bumped most dependencies

## [v0.2.1] - 2021-06-01

### Added

- Support for stdin input when ciphering
- Added new tests
- Release `snapcraft` packages
- Release `ghcr.io` images
- Publish "edge" artifacts (docker, snapcraft and binaries) for every `main` commit
- `arm64` container images

### Changed

- Replaced Drone CI with GitHub actions
- Throw a warning instead of exiting if `IPC_LOCK` is unsuccessful
- Updated to go `1.16`
- Refactored codebase following golang best practices
- Updated all dependencies

## [0.2.0] - 2020-09-28

### Added

- **Support for KV v2**
- Support of `~/.vault-token` file for Vault authentication
- gosec tests as part of the linting process
- Lock process memory before proceeding to operations with Vault API

### Changed

- Moved **logger** definition into its own package
- Moved **cli** definition into its own package
- Refactored client instanciations
- Bumped to `yaml.v3`
- Bumped to go `1.15` and goreleaser `0.143.0`
- Refactored the rand function with crypto/rand and base4
- Fixed newly discovered lint issues
- Outsourced logger configuration
- Upgraded urfave/cli to v2
- Bumped Vault to `1.5.4`
- Switched default branch to **main**
- Moved `get-secret-path` and `set-secret-path` functions under `secret get-path/set-path`
- Removed redundant config path data in statefile
- Use s5 + Vault engine as ciphering/deciphering mechanism for the local state

### Deleted

- Removed **config** package

## [0.1.8] - 2019-07-18

### Added

- `homebrew` package release
- `deb` package release
- `rpm` package release
- `scoop` package release
- `freebsd` packages

### Changed

- Fixed goimports test not breaking on errors
- Bumped Vault to **1.1.3**
- Updated go dependencies to their latest versions (2019-07-18)

### Removed

- Replaced `gox` with `goreleaser`

## [0.1.7] - 2019-03-31

### Added

- Release binaries are now automatically built and published from the CI

### Changed

- Optimized Makefile
- Upgraded Vault in test container to `1.1.0`
- Upgraded dependencies
- Fixed test coverage reports
- Moved CI from `Travis` to `Drone`

## [0.1.6] - 2019-03-28

### Added

- Also build for **arm64**

### Changed

- Fixed Dockerfile build
- Fixed Travis CI builds
- Wait a bit longer for Vault container to be ready in dev-env
- Removed the IPC_LOCK capability over the build container
- Fixed the ldflags breaking darwin and windows builds
- ignore dist folder in git
- Do not use go mod for build dependencies
- Tidied `go.mod`

## [0.1.5] - 2019-03-18

### Added

- Added gox and ghr features to release binaries

### Changed

- Fixed a panic issue on `status` and `plan` command when the Vault path doesn't contain any value
- Updated Travis CI configuration
- Upgraded Vault to `1.0.3`
- Upgraded to golang `1.12`
- Switched to `gomodules`
- Enhanced makefile
- Updated all dependencies to their latest versions
- Made the secondary container in dev-env use the same version of Vault
- Added `IPC_LOCK` capabilities to the dev-env docker container
- Upgraded Vault libraries to `0.9.6`
- Updated license to `Apache 2.0`

## [0.1.4] - 2018-02-01

### Added

- Added a flag to pass sensitive content through stdin - #8
- New function `strongbox transit delete <transit_key_name>`

### Changed

- Lint CI job was failing issue since last commits
- Fixed a bug while returning an empty transit key list
- Updated dependencies
- Support Vault `0.9.3` for development env

## [0.1.3] - 2018-01-15

### Added

- Embedded authentication against Vault using [approle](https://www.vaultproject.io/docs/auth/approle.html) auth backend - #6
- Switched base release container from empty (`scratch`) to `busybox` in order to be able to use it natively with GitLab CI\

## [0.1.2] - 2018-01-14

### Added

- Implemented a function to rotate transit keys - #2
- Got a full test environment in Makefile, added Vault container
- Possibility to generate random passwords on secret writes - #4
- Added links in Changelog

### Changed

- Nicer version output
- Updated CLI, added some flags on secret write and read functions
- Enhanced functions usage outputs
- Fixed `status` command on empty Vault Cluster

## [0.1.1] - 2018-01-13

### Added

- Added CHANGELOG.md
- Updated dependencies

### Changed

- Fixed Dockerfile
- Fixed build versioning

## [0.1.0] - 2018-01-12

### Added

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

[Unreleased]: https://github.com/mvisonneau/strongbox/compare/v0.2.2...HEAD
[v0.2.2]: https://github.com/mvisonneau/strongbox/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/mvisonneau/strongbox/compare/0.2.0...v0.2.1
[0.2.0]: https://github.com/mvisonneau/strongbox/compare/0.1.8...0.2.0
[0.1.8]: https://github.com/mvisonneau/strongbox/compare/0.1.7...0.1.8
[0.1.7]: https://github.com/mvisonneau/strongbox/compare/0.1.6...0.1.7
[0.1.6]: https://github.com/mvisonneau/strongbox/compare/0.1.5...0.1.6
[0.1.5]: https://github.com/mvisonneau/strongbox/compare/0.1.4...0.1.5
[0.1.4]: https://github.com/mvisonneau/strongbox/compare/0.1.3...0.1.4
[0.1.3]: https://github.com/mvisonneau/strongbox/compare/0.1.2...0.1.3
[0.1.2]: https://github.com/mvisonneau/strongbox/compare/0.1.1...0.1.2
[0.1.1]: https://github.com/mvisonneau/strongbox/compare/0.1.0...0.1.1
[0.1.0]: https://github.com/mvisonneau/strongbox/tree/0.1.0
