package config

import (
	"github.com/DarthSim/godotenv"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
)

type Option func(*configOptions)

type configOptions struct {
	prefix    string
	skipEnv   bool
	skipValid bool
}

func WithPrefix(prefix string) Option {
	return func(o *configOptions) {
		o.prefix = prefix
	}
}

func SkipEnvFile() Option {
	return func(o *configOptions) {
		o.skipEnv = true
	}
}

func SkipValidation() Option {
	return func(o *configOptions) {
		o.skipValid = true
	}
}

func Load(dest interface{}, opts ...Option) error {
	options := &configOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if !options.skipEnv {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	err := env.ParseWithOptions(dest, env.Options{
		Prefix: options.prefix,
	})
	if err != nil {
		return err
	}

	defaults.SetDefaults(dest)

	if !options.skipValid {
		validate := validator.New()
		err = validate.Struct(dest)
		if err != nil {
			return err
		}
	}

	return nil
}
