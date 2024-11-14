package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
	"github.com/spf13/viper"
	"news-weeder/internal/server"
	"news-weeder/internal/weeder/redis"
)

type Config struct {
	Server server.Config
	Redis  redis.RedisConfig
}

func FromFile(filePath string) (*Config, error) {
	config := &Config{}

	viperInstance := viper.New()
	viperInstance.AutomaticEnv()
	viperInstance.SetConfigFile(filePath)

	viperInstance.SetDefault("redis.Address", "localhost:6379")
	viperInstance.SetDefault("redis.Index", "news-vec-index")
	viperInstance.SetDefault("redis.KNN", 5)
	viperInstance.SetDefault("redis.DIM", 768)

	viperInstance.SetDefault("server.Address", "0.0.0.0:2866")
	viperInstance.SetDefault("server.LoggerLevel", "INFO")

	if err := viperInstance.ReadInConfig(); err != nil {
		confErr := fmt.Errorf("failed while reading config file %s: %w", filePath, err)
		return config, confErr
	}

	if err := viperInstance.Unmarshal(config); err != nil {
		confErr := fmt.Errorf("failed while unmarshaling config file %s: %w", filePath, err)
		return config, confErr
	}

	return config, nil
}

func LoadEnv(enableDotenv bool) (*Config, error) {
	if enableDotenv {
		_ = godotenv.Load()
	}

	redisAddr := loadString("NEWS_WEEDER_REDIS_ADDRESS")
	redisIndex := loadString("NEWS_WEEDER_REDIS_INDEX")
	redisKNN := loadNumber("NEWS_WEEDER_REDIS_KNN", 10)
	redisDIM := loadNumber("NEWS_WEEDER_REDIS_DIM", 10)
	redisConfig := redis.RedisConfig{
		Address: redisAddr,
		Index:   redisIndex,
		KNN:     redisKNN,
		DIM:     redisDIM,
	}

	servAddr := loadString("NEWS_WEEDER_SERVER_ADDRESS")
	servLogger := loadString("NEWS_WEEDER_SERVER_LOGGER_LEVEL")
	serverConfig := server.Config{Address: servAddr, LoggerLevel: servLogger}

	return &Config{
		Server: serverConfig,
		Redis:  redisConfig,
	}, nil
}

func loadString(envName string) string {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("faile to extract %s env var: %s", envName, value)
		log.Println(msg)
		return ""
	}
	return value
}

func loadNumber(envName string, bitSize int) int {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("faile to extract %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	number, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		msg := fmt.Sprintf("faile to convert %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	return int(number)
}

func loadBool(envName string) bool {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("faile to extract %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		msg := fmt.Sprintf("faile to convert %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	return boolean
}
