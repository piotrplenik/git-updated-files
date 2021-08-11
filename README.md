[![Build](https://github.com/jupeter/git-helpers/actions/workflows/release.yml/badge.svg)](https://github.com/jupeter/git-helpers)

# GIT Helpers

GIT Helpers help to find list of changed files between two branches.
The tool is dedicated mostly for CI/CD workflows.

The quickest way to try the command-line interface is an in-lined configuration.
```bash
# Download the latest release as /usr/local/bin/git-updated-files
$ curl https://raw.githubusercontent.com/jupeter/git-helpers/main/install.sh | bash -s -- -b /usr/local/bin
# Run the command
$ git-updated-files --target-ref master --filter ".(hcl|txt)"
```

Run `git-updated-files --help` or have a look at the [Usage Docs](USAGE.md) for more information.
