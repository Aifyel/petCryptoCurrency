package config

import (
	"time"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL       string        `mapstructure:"database_url"`
	CoingeckoURL      string        `mapstructure:"coingecko_url"`
	CoingeckoInterval time.Duration `mapstructure:"coingecko_interval"`
	GRPCPort          string        `mapstructure:"grpc_port"`
	KafkaBrokers      []string      `mapstructure:"kafka_brokers"`
	KafkaGroupID      string        `mapstructure:"kafka_group_id"`
	GRPCServerAddress string        `mapstructure:"grpc_server_address"`
	HTTPPort          string        `mapstructure:"http_port"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	bindings := map[string]string{
		"database_url":        "DATABASE_URL",
		"coingecko_url":       "COINGECKO_URL",
		"coingecko_interval":  "COINGECKO_INTERVAL",
		"grpc_port":           "GRPC_PORT",
		"kafka_brokers":       "KAFKA_BROKERS",
		"kafka_group_id":      "KAFKA_GROUP_ID",
		"grpc_server_address": "GRPC_SERVER_ADDRESS",
		"http_port":           "HTTP_PORT",
	}

	for key, envVar := range bindings {
		if err := viper.BindEnv(key, envVar); err != nil {
			return nil, errors.Wrap(entities.ErrInvalidParams, "config: "+key)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "config: unmarshal")
	}

	return &config, nil
}
