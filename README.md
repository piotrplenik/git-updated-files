[![Build](https://github.com/jupeter/git-updated-files/actions/workflows/release.yml/badge.svg)](https://github.com/jupeter/git-updated-files)
[![Coverage](https://codecov.io/gh/jupeter/git-updated-files/branch/main/graph/badge.svg)](https://codecov.io/gh/jupeter/git-updated-files)
[![Go Report Card](https://goreportcard.com/badge/github.com/jupeter/git-updated-files)](https://goreportcard.com/report/github.com/jupeter/git-updated-files)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

# GIT Helpers

GIT Helpers help to find list of changed files between two branches.
The tool is dedicated mostly for CI/CD workflows.

The quickest way to try the command-line interface is an in-lined configuration.
```bash
# Download the latest release as /usr/local/bin/git-updated-files
$ curl https://raw.githubusercontent.com/jupeter/git-updated-files/main/install.sh \
    | bash -s -- -b /usr/local/bin
# Run the command
$ git-updated-files --target-ref master --filter ".(hcl|txt)"
```

Run `git-updated-files --help` or have a look at the [Usage Docs](USAGE.md) for more information.
