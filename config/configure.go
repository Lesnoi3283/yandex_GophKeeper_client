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
	defaultAPIAddress          = "localhost:8080"
	defaultLogLevel            = "info"
	defaultMaxBinDataChunkSize = 16
	defaultUseHTTPS            = true
	defaultGRPCAddress         = "localhost:50051"
)

type AppConfig struct {
	APIAddress          string
	LogLevel            string
	MaxBinDataChunkSize int
	UseHTTPS            bool
	GRPCAddress         string
}

// Configure reads configuration params from environment variables and command line arguments.
//
//	Priority (1 - high, N - low):
//	1 - Environment.
//	2 - flags.
func (c *AppConfig) Configure() error {
	// Get flag values
	flag.StringVar(&(c.APIAddress), "api-address", defaultAPIAddress, "API server address without protocol. Example: \"localhost:8080\" (not a \"https://localhost:8080\".")
	flag.StringVar(&(c.LogLevel), "log-level", defaultLogLevel, "Log level.")
	flag.IntVar(&(c.MaxBinDataChunkSize), "max-bin-data-chunk-size", defaultMaxBinDataChunkSize, "Max size of binary data chunks in bytes.")
	flag.BoolVar(&(c.UseHTTPS), "use-https", defaultUseHTTPS, "Use HTTPS.")
	flag.StringVar(&(c.GRPCAddress), "grpc-address", defaultGRPCAddress, "full GRPC server address with port.")
	flag.Parse()

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
	if value, found := os.LookupEnv("USE_HTTPS"); found {
		valueBool, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("error parsing USE_HTTPS: %w", err)
		}
		c.UseHTTPS = valueBool
	}
	if value, found := os.LookupEnv("GRPC_ADDRESS"); found {
		c.GRPCAddress = value
	}

	return nil
}
