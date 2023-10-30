package config

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	App      ApplicationConfig
	Database DatabaseConfig
}

func NewConfig() *Config {
	config := &Config{}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetDefault("config-path", "./config.yaml")
	configPath := viper.GetString("config-path")
	if _, err := os.Stat(configPath); err == nil {
		viper.SetConfigType(filepath.Ext(configPath)[1:])
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("failed to read config (%s): %+v", configPath, err)
		}
	}

	config.Database.parse()
	config.App.parse()

	return config
}

type ApplicationConfig struct {
	ServiceName    string
	WebPort        uint
	HttpPathPrefix string

	LogLevel string
}

func (c *ApplicationConfig) parse() {
	c.ServiceName = viperGetOrDefault("app.service-name", "OHA")
	c.WebPort = viperGetOrDefaultUint("app.web-port", 8000)
	c.HttpPathPrefix = viperGetOrDefault("app.http-path-prefix", "/api/v1")
	c.LogLevel = viperGetOrDefault("app.log-level", "debug")
}

type DatabaseConfig struct {
	Host     string
	Port     uint
	User     string
	Password string
	DB       string
}

func (c *DatabaseConfig) parse() {
	c.Host = viperGetOrDefault("database.host", "127.0.0.1")
	c.Port = viperGetOrDefaultUint("database.port", 5432)
	c.User = viperGetOrDefault("database.user", "postgres")
	c.Password = viperGetOrDefault("database.password", "12345")
	c.DB = viperGetOrDefault("database.db", "oha")
}

func viperGetOrDefault(key string, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}

func viperGetOrDefaultUint(key string, defaultValue uint64) uint {
	viper.SetDefault(key, defaultValue)
	return viper.GetUint(key)
}
