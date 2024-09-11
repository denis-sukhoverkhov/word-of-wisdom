package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Difficulty      uint8
	Addr            string
	WorkerCount     int
	ShutdownTimeout int
	PowTimeout      int
}

type ClientConfig struct {
	ServerAddr    string
	RPS           int
	TotalRequests int
}

func loadConfig(configName, configDir, configType string) error {
	viper.SetConfigType(configType)
	viper.AddConfigPath(configDir)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// check if env is set: ENV ("prod", "stage", "test")
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	viper.SetConfigName(fmt.Sprintf("%s_%s", configName, env))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}

func LoadServerConfig(configDir string) (*ServerConfig, error) {
	if err := loadConfig("server_config", configDir, "yaml"); err != nil {
		return nil, err
	}

	config := &ServerConfig{
		Difficulty:      uint8(viper.GetInt("DIFFICULTY")),
		Addr:            viper.GetString("ADDR"),
		WorkerCount:     viper.GetInt("WORKERCOUNT"),
		ShutdownTimeout: viper.GetInt("SHUTDOWNTIMEOUT"),
		PowTimeout:      viper.GetInt("POWTIMEOUT"),
	}

	return config, nil
}

func LoadClientConfig(configDir string) (*ClientConfig, error) {
	if err := loadConfig("client_config", configDir, "yaml"); err != nil {
		return nil, err
	}

	config := &ClientConfig{
		ServerAddr:    viper.GetString("SERVERADDR"),
		RPS:           viper.GetInt("RPS"),
		TotalRequests: viper.GetInt("TOTALREQUESTS"),
	}

	return config, nil
}
