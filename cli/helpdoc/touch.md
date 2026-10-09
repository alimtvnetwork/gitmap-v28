# touch

Create an empty file or update file modification timestamps, automatically creating missing parent directories.

## Aliases

`mkfile`, `create-file`

## Usage

    gitmap touch <filepath>
    gitmap mkfile <filepath> [content...]

## Description

The `touch` command creates empty files or updates their timestamp across all operating systems. If parent directories do not exist, GitMap automatically provisions them.

`mkfile` also allows writing initial content directly to the file upon creation.

## Examples

```bash
# Create a new empty file with parent directories
gitmap touch src/config/local.yaml

# Create a file with initial content
gitmap mkfile configs/app.env "APP_ENV=dev PORT=8080"
```
