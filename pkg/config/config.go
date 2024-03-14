package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	devEnv          = "production"
	envConfigPrefix = "CONFIG"
)

// Config gets and sets values from a config file
type Config struct {
	viper *viper.Viper
}

// Load and read from configuration with a default environment.
func Load(fileName string, paths []string) (Config, error) {
	vi := viper.New()

	// Setting the prefix to load from ENV
	vi.SetEnvPrefix(envConfigPrefix)
	// This will cause viper to automatically check for the key prefix_KEY and return value
	vi.AutomaticEnv()

	cfg := Config{vi}
	if fileName == "" {
		return cfg, fmt.Errorf("config: no env specified")
	}
	cfg.viper.SetConfigName(fileName)
	for _, path := range paths {
		cfg.viper.AddConfigPath(path)
	}
	if err := cfg.viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("config read file %s failed: %s", fileName, err)
	}
	return cfg, nil
}

// Env retreives the current environment based on the ACTIVE_ENV variable
func Env() string {
	if env := os.Getenv("ACTIVE_ENV"); env != "" {
		return env
	}
	return devEnv
}

// GetInt fetches the value with the afformentioned type
func (cfg Config) GetInt(key string) int {
	return cfg.viper.GetInt(key)
}

// GetString fetches the value with the afformentioned type
func (cfg Config) GetString(key string) string {
	return cfg.viper.GetString(key)
}

// GetStringSlice fetches the value with the afformentioned type
func (cfg Config) GetStringSlice(key string) []string {
	return cfg.viper.GetStringSlice(key)
}
