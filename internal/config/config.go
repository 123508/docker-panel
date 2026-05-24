package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	User   UserConfig   `mapstructure:"user"`
	Server ServerConfig `mapstructure:"server"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type UserConfig struct {
	AdminUsername string `mapstructure:"admin_username"`
	AdminPassword string `mapstructure:"admin_password"`
}

type ServerConfig struct {
	BindIP   string `mapstructure:"bind_ip"`
	BindPort string `mapstructure:"bind_port"`
	Debug    bool   `mapstructure:"debug"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
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

func SetAdminPassword(password string) error {
	viper.Set("user.admin_password", password)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if AppConfig != nil {
		AppConfig.User.AdminPassword = password
	}

	return nil
}

func createDefaultConfig(configFile string) error {
	configPath := filepath.Dir(configFile)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	adminPassword, err := generateRandomPassword()
	if err != nil {
		return err
	}

	jwtSecret, err := generateJWTSecret()
	if err != nil {
		return err
	}

	defaultConfig := fmt.Sprintf(`[user]
admin_username = "admin"
admin_password = "%s"

[server]
bind_ip = "0.0.0.0"
bind_port = "8080"
debug = false

[jwt]
secret = "%s"
`, adminPassword, jwtSecret)

	return os.WriteFile(configFile, []byte(defaultConfig), 0644)
}

func generateRandomPassword() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func generateJWTSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate JWT secret: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func getConfigPath() string {
	if cwd, err := os.Getwd(); err == nil {
		if fileExists(filepath.Join(cwd, "config.toml")) {
			return cwd
		}
	}

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
