# Config Package

A flexible configuration package for Go applications that combines environment variables, `.env` files, default values, and validation.

## Features

- Load configuration from environment variables
- Support for `.env` files
- Set default values using struct tags
- Validate configuration using struct tags
- Configurable prefix for environment variables
- Optional validation and `.env` file loading

## Installation

```bash
go get github.com/kfreiman/x/config
```

## Usage

### Basic Usage

```go
type Config struct {
    Host     string `env:"HOST" validate:"required"`
    Port     int    `env:"PORT" default:"8080"`
    Debug    bool   `env:"DEBUG" default:"false"`
    Optional string `env:"OPTIONAL"`
}

func main() {
    cfg := &Config{}
    err := config.Load(cfg)
    if err != nil {
        log.Fatal(err)
    }
}
```

### With Options

```go
cfg := &Config{}
err := config.Load(cfg,
    config.WithPrefix("APP_"),    // Use APP_ prefix for env vars
    config.SkipEnvFile(),        // Skip loading .env file
    config.SkipValidation(),     // Skip validation
)
```

### Environment Variables

Set environment variables directly or via `.env` file:

```env
HOST=localhost
PORT=9090
DEBUG=true
```

## Options

- `WithPrefix(prefix string)`: Add prefix to environment variable names
- `SkipEnvFile()`: Skip loading `.env` file
- `SkipValidation()`: Skip struct validation

## Dependencies

- [godotenv](https://github.com/DarthSim/godotenv) - For `.env` file support
- [env](https://github.com/caarlos0/env) - For environment variable parsing
- [validator](https://github.com/go-playground/validator) - For struct validation
- [defaults](https://github.com/mcuadros/go-defaults) - For default values

## License

MIT License
