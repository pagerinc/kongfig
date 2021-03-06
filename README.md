# kongfig

Kongfig is a configuration management tool for the Kong API gateway.

## Usage

```bash
kongfig apply -f examples/config.example.json --dry-run
```

### Available Commands

| Command   | Description                              |
| ---       | ---                                      |
| `apply`   | Apply a configuration to a Kong instance |
| `help`    | Help about any command                   |
| `version` | Print the version number of Kongfig      |

Use `kongfig [command] --help` for more information about a command.

## Contributing

1. Fork the project
2. Make your changes
3. `make`
4. The `kongfig` binary is now available to use locally

> **NOTE**: when building in OS X, you'll need to export the `GOOS` env, eg:

```bash
GOOS=darwin make
```

Additionaly, there's a docker-compose file included for ease of local
development.
Simply call `docker-compose up` to get started.

Dependencies are managed using [dep]. Please refer to its documentation if needed.

### Testing

Tests are run automatically on every build or via the `make test` target.
Additionaly you can run `make cover` to check your coverage.

[dep]: https://github.com/golang/dep
