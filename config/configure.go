// Package config is a configuration package.
// It contains a AppConfig struct and functions to read configuration params from env variables and command line args.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Default configuration values.
const (
	DefaultAPIAddress          = "localhost:8080"
	DefaultLogLevel            = "info"
	DefaultMaxBinDataChunkSize = 16
)

type AppConfig struct {
	UserDataPath        string
	APIAddress          string
	LogLevel            string
	MaxBinDataChunkSize int
}

// Configure reads configuration params from environment variables and command line arguments.
func (c *AppConfig) Configure() error {
	// Get flag values
	flag.StringVar(&(c.UserDataPath), "user-data-path", "", "Path to user data directory (required).")
	flag.StringVar(&(c.APIAddress), "api-address", DefaultAPIAddress, "API server address. Example: \"localhost:8080\".")
	flag.StringVar(&(c.LogLevel), "log-level", DefaultLogLevel, "Log level.")
	flag.IntVar(&(c.MaxBinDataChunkSize), "max-bin-data-chunk-size", DefaultMaxBinDataChunkSize, "Max size of binary data chunks in bytes.")
	flag.Parse()

	// Get env values
	if value, found := os.LookupEnv("USER_DATA_PATH"); found {
		c.UserDataPath = value
	}
	if value, found := os.LookupEnv("API_ADDRESS"); found {
		c.APIAddress = value
	}
	if value, found := os.LookupEnv("LOG_LEVEL"); found {
		c.LogLevel = value
	}
	if value, found := os.LookupEnv("MAX_BIN_DATA_CHUNK_SIZE"); found {
		size, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("error parsing MAX_BIN_DATA_CHUNK_SIZE: %w", err)
		}
		c.MaxBinDataChunkSize = size
	}

	// Validate required fields
	if c.UserDataPath == "" {
		return fmt.Errorf("USER_DATA_PATH is required and must be set either as an environment variable or command line argument")
	}

	return nil
}
