# bulker

Runs different operations on a bunch of repositories in bulk mode

[![Release](https://img.shields.io/github/v/release/mih-kopylov/bulker?style=for-the-badge)](https://github.com/mih-kopylov/bulker/releases/latest)
[![GitHub license](https://img.shields.io/github/license/mih-kopylov/bulker?style=for-the-badge)](https://github.com/mih-kopylov/bulker/blob/master/LICENSE)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/mih-kopylov/bulker/build?style=for-the-badge)](https://github.com/mih-kopylov/bulker/actions/workflows/build.yml)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/mih-kopylov/bulker)

## Quick Start

### Install

If GO installed, run

```shell
go install github.com/mih-kopylov/bulker/app/bulker@latest
```

Or just download an appropriate binary from [Assets](https://github.com/mih-kopylov/bulker/releases/latest)

### Configure

Bulker works without any preliminary configuration.

Just create a directory for the repositories, say `~/projects` and `cd` into it. Bulker stores repositories inside
current directory by default.

### Run

```shell
bulker repos add --name bulker --url https://github.com/mih-kopylov/bulker --tags "github,test"
bulker repos list -t github
bulker git clone
bulker status
```