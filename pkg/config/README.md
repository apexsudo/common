# config

Generic config loader that reads JSON, TOML, and YAML files into typed Go structs, with environment variable override support via [configor](https://github.com/jinzhu/configor).

## Usage

Define a struct that mirrors your config file, then call `LoadConfig`:

```go
import "github.com/apexsudo/common/pkg/config"

type AppConfig struct {
    Port int    `env:"APP_PORT" json:"port" default:"5432"`
    DSN  string `env:"DB_DSN" json:"dsn" required:"true"`
}

cfg, err := config.LoadConfig[AppConfig]("app.json")
if err != nil {
    log.Fatal(err)
}
```

Multiple files are merged in order, with later files taking precedence:

```go
cfg, err := config.LoadConfig[AppConfig]("base.yaml", "overrides.yaml")
```

## Config directory resolution

Files are resolved relative to a base directory determined at runtime:

1. **`CONFIG_PATH` env var** — if set, files are loaded from that directory.
2. **Source file directory** — fallback to the directory containing this package's source, useful when config files are bundled alongside the source.

## Environment variable overrides

Any field can be overridden at runtime via environment variables. configor maps struct fields to env vars. See the [configor docs](https://github.com/jinzhu/configor) for full details.
