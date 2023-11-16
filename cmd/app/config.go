package main

import (
	"fmt"
	"strings"

	"github.com/fedev521/go-net-gen/internal/app"
	"github.com/fedev521/go-net-gen/internal/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Overall program configuration.
type configuration struct {
	Log log.Config
	App app.Config
}

// Validate validates the configuration.
func (c configuration) Validate() error {

	if err := c.App.Validate(); err != nil {
		return err
	}

	return nil
}

// Process post-processes configuration after loading it.
func (c configuration) Process() error {
	return nil
}

// Configures default settings of configuration and flags.
func configureDefaultSettings() {
	// config file settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("configs")
	viper.AddConfigPath("../../configs")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// environment variable settings
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	// flags configuration
	pflag.String("hub-project-id", "", "Hub project id on GCP")
	viper.BindPFlag("app.hub_project_id", pflag.Lookup("hub-project-id"))

	// config defaults
	viper.SetDefault("log.level", "info")
}

// Reads, unmarshals and post-processes configuration.
func loadConfiguration() (configuration, error) {
	// read configuration
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return configuration{}, fmt.Errorf("config file not found: %w", err)
		} else {
			return configuration{}, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// unmarshal configuration
	var config configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return configuration{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// post-process configuration
	err = config.Process()
	if err != nil {
		return configuration{}, fmt.Errorf("failed to post-process configuration: %w", err)
	}

	return config, nil
}
