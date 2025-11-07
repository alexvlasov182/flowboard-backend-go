package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
	Mode      string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{
		Port:      viper.GetString("PORT"),
		DBHost:    viper.GetString("DB_HOST"),
		DBPort:    viper.GetString("DB_PORT"),
		DBUser:    viper.GetString("DB_USER"),
		DBPass:    viper.GetString("DB_PASS"),
		DBName:    viper.GetString("DB_NAME"),
		JWTSecret: viper.GetString("JWT_SECRET"),
		Mode:      viper.GetString("GIN_MODE"),
	}
	return cfg, nil
}
