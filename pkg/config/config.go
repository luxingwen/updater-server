package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	LogConfig  LogConfig
	MySQL      MySQLConfig
}

type LogConfig struct {
	Level    string
	Format   string
	File     string // 日志文件路径
	MaxSize  int    // 最大文件大小（MB）
	MaxAge   int    // 最大文件保留天数
	Compress bool   // 是否压缩
	Filename string
}

type MySQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

var (
	config *Config
)

func InitConfig() {
	loadConfigFile()
	bindEnvs()
	unmarshalConfig()
}

func loadConfigFile() {
	v := viper.New()

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}

	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
}

func bindEnvs() {
	viper.AutomaticEnv()

	viper.BindEnv("ServerPort", "SERVER_PORT")
	viper.BindEnv("LogConfig.Level", "LOG_LEVEL")
	viper.BindEnv("LogConfig.Format", "LOG_FORMAT")
	viper.BindEnv("LogConfig.File", "LOG_FILE")
	viper.BindEnv("LogConfig.MaxSize", "LOG_MAX_SIZE")
	viper.BindEnv("LogConfig.MaxAge", "LOG_MAX_AGE")
	viper.BindEnv("LogConfig.Compress", "LOG_COMPRESS")
	viper.BindEnv("MySQL.Host", "MYSQL_HOST")
	viper.BindEnv("MySQL.Port", "MYSQL_PORT")
}

func unmarshalConfig() {
	config = &Config{}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
}

func GetConfig() *Config {
	return config
}
