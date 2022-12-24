# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Getting defualt remote branch
- Getting branches when repository is in detached HEAD state
- Do not clear locally stored repositories if import fails

## [0.10.0] - 2022-12-05

### Added

- `force` flag for `git push` command
- Number of errors indication in the progress bar
- `from` flag for `group create` command
- `previous` group is created with the repositories of the previous command run
- `configure` command

### Changed

- `repo` to `name` flag name to mention repository names in `groups create`, `groups append`, `groups exclude` commands

## [0.9.0] - 2022-11-04

### Added

- `open` command
- `files` commands: `copy`, `rename`, `remove`, `search`, `replace`

## [0.8.0] - 2022-08-06

### Added

- `git push` command
- `git commit` command

### Fixed

- Empty array to be printed in `--output json` mode when no repositories matched

## [0.7.0] - 2022-08-05

### Added

- `gocyclo` to the build pipeline

### Changed

- Commands that process repositories are terminated gracefully

## [0.6.0] - 2022-08-02

### Added

- `git fetch` command
- `git pull` command
- Singular aliases for `groups` and `repos` commands
- `git branches` commands: `list`, `checkout`, `create`, `remove`, `clean`
- `--ref` flag for `status` command

### Changed

- `OK` status to `Clean`
- Replace `--ok`, `--dirty` and `--missing` flags with `--show` in `status` command
- Negate prefix for filter from `-` to `!`

### Fixed

- `--name` filter regexp to match full repository name only

## [0.5.0] - 2022-07-29

### Added

- `table` output mode
- Filter repositories by group

### Changed

- `--name` parameter for repositories filtering to consume regexp

## [0.4.0] - 2022-07-29

### Added

- `--recreate` flag for `git clone` command that recreates previously cloned repository
- `status` command
- `run` command

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
