# unveil

[![Go Report Card](https://goreportcard.com/badge/github.com/gi8lino/unveil?style=flat-square)](https://goreportcard.com/report/github.com/gi8lino/unveil)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gi8lino/unveil)
[![Release](https://img.shields.io/github/release/gi8lino/unveil.svg?style=flat-square)](https://github.com/gi8lino/unveil/releases/latest)
![Tests](https://github.com/gi8lino/unveil/actions/workflows/tests.yml/badge.svg)
[![Build](https://github.com/gi8lino/unveil/actions/workflows/release.yml/badge.svg)](https://github.com/gi8lino/unveil/actions/workflows/release.yml)
[![license](https://img.shields.io/github/license/gi8lino/unveil.svg?style=flat-square)](LICENSE)

---

`unveil` is a command-line tool to extract values from configuration files (`JSON`, `YAML`, `TOML`, `INI`, `.env`/key-value files) and expose them as environment variables.
It is designed for scripting and container environments where you want to _unveil_ secrets or config values and pass them to processes in a standardized `KEY=VALUE` format.

## Features

- Supports multiple file formats:
  - **JSON**
  - **YAML**
  - **TOML**
  - **INI**
  - **key=value files** (including `.env` with `export KEY=VAL`)
- Flexible key selectors:
  - Dot notation for nested fields: `server.host`
  - Array index: `servers.0.host`
  - Array filter: `servers.[name=db].port`
- Output options:
  - Plain `KEY=VALUE`
  - With `export` prefix (`--export`)
  - Optional quoting (`none`, `single`, `double`, `json`)
- Atomic file output with `--output` (safe for CI/CD)
- Multiple extractions in one command (dynamic groups)

## Installation

```bash
go install github.com/gi8lino/unveil@latest
```

or download binaries from [Releases](https://github.com/gi8lino/unveil/releases).

## Usage

```bash
unveil [flags]
```

### Examples

Extract a single value from JSON:

```bash
unveil \
  --json.db.path=./config.json \
  --json.db.select=database.user \
  --json.db.as=DBUSER
```

Output:

```bash
DBUSER=alice
```

Extract from YAML and TOML, with quoting:

```bash
unveil \
  --yaml.srv.path=./app.yaml \
  --yaml.srv.select=server.host \
  --yaml.srv.as=HOST \
  --yaml.srv.quote=single \
  --toml.t.path=./cfg.toml \
  --toml.t.select=server.port \
  --toml.t.as=PORT
```

Output:

```bash
HOST='localhost'
PORT=8080
```

Write directly to a file instead of stdout:

```bash
unveil --file.env.path=.env --file.env.select=TOKEN --file.env.as=API_TOKEN --output=out.env
```

Result (`out.env`):

```bash
API_TOKEN=secret
```

With `--export`:

```bash
export API_TOKEN=secret
```

### Supported Flags

- `--quote MODE` — global quote mode for all values
  One of: `none`, `single`, `double`, `json`
- `--output FILE` — write results atomically to `FILE` instead of stdout
- `--export` — prefix each line with `export `

Each extractor group (`json`, `yaml`, `toml`, `ini`, `file`) supports:

- `--<group>.<id>.path=PATH` (required)
- `--<group>.<id>.select=KEY` (required)
- `--<group>.<id>.as=VAR` (optional, defaults to uppercase ID)
- `--<group>.<id>.quote=MODE` (optional override)

## Development

Run unit tests:

```bash
go test ./...
```

Lint:

```bash
golangci-lint run
```

## License

Apache 2.0 -- see [LICENSE](LICENSE)
