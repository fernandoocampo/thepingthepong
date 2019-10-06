package domain

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	// ConfigurationFileName contains configuration file name
	ConfigurationFileName = "config"
)

// Configuration contains the configuration data.
var Configuration Setting

// LogData contains data configuration for logs
type LogData struct {
	Level  string
	Format string
}

// LogSetting contains configuration for log
type LogSetting struct {
	Main       LogData // Log configuration for Main module
	Port       LogData // Log configuration for Port module
	Domain     LogData // Log configuration for Domain module
	Authapp    LogData // Log configuration for AuthApp module
	Matchapp   LogData // Log configuration for MatchApp module
	Playerapp  LogData // Log configuration for PlayerApp module
	Repository LogData // Log configuration for Repository module
}

// ServerSetting contains the configuration parameters.
type ServerSetting struct {
	Port string // Web server port
}

// Setting contains general configuration data for the application.
type Setting struct {
	Log       LogSetting    // configuration data for log
	Webserver ServerSetting // configuration data for server
}

// LoadConfiguration creates a new configuration
func LoadConfiguration(configPath string) {
	loadConfigurationFile(configPath)
	loadConfigurationModel()
}

// loadConfigurationModel load the configuration data into a model
func loadConfigurationModel() {
	errUnmarshal := viper.Unmarshal(&Configuration)
	if errUnmarshal != nil {
		log.Fatalf("unable to decode configuration file, %v", errUnmarshal)
	}
}

// loadConfigurationFile loads the file into memory
func loadConfigurationFile(configPath string) {
	viper.SetConfigName(ConfigurationFileName)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Config file was found but another error was produced
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
