package config

import (
	"github.com/craftizmv/rewards/internal/data/infrastructure/cache"
	"github.com/craftizmv/rewards/internal/data/infrastructure/database"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/craftizmv/rewards/server"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	ServiceName string                `mapstructure:"serviceName"`
	Logger      *logger.LoggerConfig  `mapstructure:"logger"`
	Rabbitmq    *queue.RabbitMQConfig `mapstructure:"rabbitmq"`
	EchoCfg     *server.EchoConfig    `mapstructure:"echo"`
	CacheCfg    *cache.Config         `mapstructure:"cache"`
	DBCfg       *database.Config      `mapstructure:"db"`
}

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config_dev")
		viper.SetConfigType("json")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}
	})

	return configInstance
}
