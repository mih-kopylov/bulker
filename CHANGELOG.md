# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `--recreate` flag for `git clone` command that recreates previously cloned repository

## [0.3.0] - 2022-07-28

### Added

- Documentation generation in GitHub Wiki
- Ability to limit number of simultaneously processed repositories in parallel run mode
- Progress bar indicating how many repositories have been processed

### Changed

- Default output mode to `line`

## [0.2.0] - 2022-07-26

### Added

- `groups` commands: `list`, `get`, `create`, `append`, `exclude`, `remove`, `clean`

## [0.1.1] - 2022-07-25

### Fixed

- Parallel runs

## [0.1.0] - 2022-07-20

### Added

- `repos` commands: `add`, `remove`, `list`, `export`, `import`
- `git clone` command
