package main

import (
	"flag"
	"os"
	"strconv"
)

// AppConfig holds application configuration.
type AppConfig struct {
	JSONPath string
	OutPath  string
	Indent   int
	Debug    bool
}

// getAppConfig initializes AppConfig from environment variables and command-line flags.
func getAppConfig() *AppConfig {
	config := &AppConfig{}

	// Define CLI flags
	flag.StringVar(&config.JSONPath, "p", "file.json", "JSON file path (Env: INPUT_PATH)")
	flag.StringVar(&config.OutPath, "o", "", "YAML output path, if empty, output to std (Env: OUTPUT_PATH)")
	flag.IntVar(&config.Indent, "i", 2, "Indentation in YAML (Env: INDENT)")
	flag.BoolVar(&config.Debug, "d", false, "Enable debug mode (Env: DEBUG)")
	flag.Parse()

	// Override with environment variable if set
	if envVal, exists := os.LookupEnv("INPUT_PATH"); exists {
		config.JSONPath = envVal
	}
	if envVal, exists := os.LookupEnv("OUTPUT_PATH"); exists {
		config.OutPath = envVal
	}
	if envVal, exists := os.LookupEnv("INDENT"); exists {
		if val, err := strconv.Atoi(envVal); err == nil {
			config.Indent = val
		}
	}
	if envVal, exists := os.LookupEnv("DEBUG"); exists {
		config.Debug = envVal != "" && envVal != "0"
	}

	return config
}
