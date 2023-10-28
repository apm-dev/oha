package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

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
	WebPort        uint64
	HttpPathPrefix string

	DatabaseConnTimeout time.Duration

	LogLevel string
}

func (c *ApplicationConfig) parse() {
	c.ServiceName = viperGetOrDefault("app.service-name", "OHA")
	c.WebPort = viperGetOrDefaultUint64("app.web-port", 8000)
	c.HttpPathPrefix = viperGetOrDefault("app.http-path-prefix", "")
	c.DatabaseConnTimeout = viperGetOrDefaultTimeDuration("app.outgoing-request-timeout", "15s")
	c.LogLevel = viperGetOrDefault("app.log-level", "debug")
}

type DatabaseConfig struct {
	Host     string
	Port     uint64
	User     string
	Password string
	DB       string
}

func (c *DatabaseConfig) parse() {
	c.Host = viperGetOrDefault("database.host", "127.0.0.1")
	c.Port = viperGetOrDefaultUint64("database.port", 5432)
	c.User = viperGetOrDefault("database.user", "postgres")
	c.Password = viperGetOrDefault("database.password", "12345")
	c.DB = viperGetOrDefault("database.db", "oha")
}

func viperGetOrDefault(key string, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}

func viperGetOrDefaultUint64(key string, defaultValue uint64) uint64 {
	viper.SetDefault(key, defaultValue)
	return viper.GetUint64(key)
}

func viperGetOrDefaultTimeDuration(key string, defaultValue string) time.Duration {
	viper.SetDefault(key, defaultValue)
	d, err := time.ParseDuration(viper.GetString(key))
	if err != nil {
		log.Fatalf("provided value '%s' cannot be transformed to [time.Duration]", viper.GetString(key))
	}
	return d
}
