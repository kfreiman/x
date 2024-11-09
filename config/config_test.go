package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Host     string `env:"HOST" validate:"required"`
	Port     int    `env:"PORT" default:"8080"`
	Debug    bool   `env:"DEBUG" default:"false"`
	Optional string `env:"OPTIONAL"`
}

func TestLoad(t *testing.T) {
	t.Run("basic configuration with defaults", func(t *testing.T) {
		cfg := &TestConfig{}
		err := Load(cfg, SkipEnvFile())
		require.Error(t, err) // Should error due to missing required HOST
		assert.Equal(t, 8080, cfg.Port)
		assert.False(t, cfg.Debug)
		assert.Empty(t, cfg.Optional)
	})

	t.Run("with environment variables", func(t *testing.T) {
		// Setup
		os.Setenv("HOST", "localhost")
		os.Setenv("PORT", "9090")
		os.Setenv("DEBUG", "true")
		defer func() {
			os.Unsetenv("HOST")
			os.Unsetenv("PORT")
			os.Unsetenv("DEBUG")
		}()

		cfg := &TestConfig{}
		err := Load(cfg, SkipEnvFile())
		require.NoError(t, err)

		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, 9090, cfg.Port)
		assert.True(t, cfg.Debug)
	})

	t.Run("with prefix", func(t *testing.T) {
		// Setup
		os.Setenv("APP_HOST", "localhost")
		os.Setenv("APP_PORT", "9090")
		defer func() {
			os.Unsetenv("APP_HOST")
			os.Unsetenv("APP_PORT")
		}()

		cfg := &TestConfig{}
		err := Load(cfg, SkipEnvFile(), WithPrefix("APP_"))
		require.NoError(t, err)

		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, 9090, cfg.Port)
	})

	t.Run("validation failure", func(t *testing.T) {
		cfg := &TestConfig{}
		err := Load(cfg, SkipEnvFile())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation")
	})

	t.Run("skip validation", func(t *testing.T) {
		cfg := &TestConfig{}
		err := Load(cfg, SkipEnvFile(), SkipValidation())
		require.NoError(t, err)
	})

	t.Run("multiple options", func(t *testing.T) {
		// Setup
		os.Setenv("APP_HOST", "localhost")
		os.Setenv("APP_PORT", "9090")
		defer func() {
			os.Unsetenv("APP_HOST")
			os.Unsetenv("APP_PORT")
		}()

		cfg := &TestConfig{}
		err := Load(cfg,
			SkipEnvFile(),
			WithPrefix("APP_"),
			SkipValidation(),
		)
		require.NoError(t, err)

		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, 9090, cfg.Port)
	})
}

func TestConfigOptions(t *testing.T) {
	t.Run("WithPrefix", func(t *testing.T) {
		opts := &configOptions{}
		WithPrefix("TEST_")(opts)
		assert.Equal(t, "TEST_", opts.prefix)
	})

	t.Run("SkipEnvFile", func(t *testing.T) {
		opts := &configOptions{}
		SkipEnvFile()(opts)
		assert.True(t, opts.skipEnv)
	})

	t.Run("SkipValidation", func(t *testing.T) {
		opts := &configOptions{}
		SkipValidation()(opts)
		assert.True(t, opts.skipValid)
	})
}

func TestEnvFile(t *testing.T) {
	t.Run("with env file", func(t *testing.T) {
		// Create temporary .env file
		envContent := []byte("HOST=testhost\nPORT=5000\n")
		err := os.WriteFile(".env", envContent, 0644)
		require.NoError(t, err)
		defer os.Remove(".env")

		cfg := &TestConfig{}
		err = Load(cfg)
		require.NoError(t, err)

		assert.Equal(t, "testhost", cfg.Host)
		assert.Equal(t, 5000, cfg.Port)
	})

	t.Run("missing env file", func(t *testing.T) {
		// Ensure .env file doesn't exist
		os.Remove(".env")

		cfg := &TestConfig{}
		err := Load(cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), ".env")
	})
}

func TestInvalidEnvironmentVariables(t *testing.T) {
	t.Run("invalid port type", func(t *testing.T) {
		os.Setenv("PORT", "invalid")
		defer os.Unsetenv("PORT")

		cfg := &TestConfig{}
		err := Load(cfg, SkipEnvFile())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parse error on field \"Port\"")
	})

	t.Run("nil destination", func(t *testing.T) {
		err := Load(nil, SkipEnvFile())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected a pointer")
	})

	t.Run("non-pointer destination", func(t *testing.T) {
		cfg := TestConfig{}
		err := Load(cfg, SkipEnvFile())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pointer")
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty prefix", func(t *testing.T) {
		cfg := &TestConfig{}
		err := Load(cfg, WithPrefix(""))
		require.Error(t, err) // Should error due to missing required HOST
		// The test passes as long as it doesn't panic
	})

	t.Run("duplicate options", func(t *testing.T) {
		cfg := &TestConfig{}
		err := Load(cfg,
			WithPrefix("APP_"),
			WithPrefix("TEST_"), // Second prefix should override first
			SkipEnvFile(),
		)
		require.Error(t, err) // Should error due to missing required HOST
		// The test passes as long as it doesn't panic
	})
}

func BenchmarkLoad(b *testing.B) {
	cfg := &TestConfig{}
	os.Setenv("HOST", "localhost")
	defer os.Unsetenv("HOST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Load(cfg, SkipEnvFile(), SkipValidation())
	}
}
