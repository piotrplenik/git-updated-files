# GIT Updated Files overview

## Available flags for `git-updated-files`

| Name | Usage | Default |
| ---- | ----- | ------- |
| `--source-ref` | GIT source reference | Current branch |
| `--target-ref` | GIT target reference | `main` |
| `--filter` | Regular expression for files matching | `.*` |
| `--directories-only` | Display directories instead of files | `true` |
| `--deleted` | Show deleted instead of added/updated | `false` |

## Sample usage

Get all `terraform` changed directories:
```
$ git-updated-files --target-ref master --filter ".hcl"
```
