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
}

func Load() (*Config, error) {
	err := viper.BindEnv("database_url", "DATABASE_URL")
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "config: database_url")
	}

	err = viper.BindEnv("coingecko_url", "COINGECKO_URL")
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "config: coingecko_url")
	}

	err = viper.BindEnv("coingecko_interval", "COINGECKO_INTERVAL")
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "config: coingecko_interval")
	}

	err = viper.BindEnv("grpc_port", "GRPC_PORT")
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "config: grpc_address")
	}

	err = viper.BindEnv("kafka_brokers", "KAFKA_BROKERS")
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "config: kafka_brokers")
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "config: unmarshal")
	}

	return &config, nil
}
