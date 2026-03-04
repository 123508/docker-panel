package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	User   UserConfig   `mapstructure:"user"`
	Server ServerConfig `mapstructure:"server"`
}

type UserConfig struct {
	AdminUsername string `mapstructure:"admin_username"`
	AdminPassword string `mapstructure:"admin_password"`
}

type ServerConfig struct {
	BindIP   string `mapstructure:"bind_ip"`
	BindPort string `mapstructure:"bind_port"`
}

var AppConfig *Config

func InitConfig() error {
	configPath := getConfigPath()
	configFile := filepath.Join(configPath, "config.toml")

	if !fileExists(configFile) {
		if err := createDefaultConfig(configFile); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
		fmt.Printf("✓ Created default config file: %s\n", configFile)
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	fmt.Printf("✓ Loaded config from: %s\n", configFile)
	return nil
}

func createDefaultConfig(configFile string) error {
	configPath := filepath.Dir(configFile)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	defaultConfig := `[user]
admin_username = "admin"
admin_password = "admin123"

[server]
bind_ip = "0.0.0.0"
bind_port = "8080"
`

	return os.WriteFile(configFile, []byte(defaultConfig), 0644)
}

func getConfigPath() string {
	if exePath, err := os.Executable(); err == nil {
		return filepath.Dir(exePath)
	}

	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}

	return "."
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
