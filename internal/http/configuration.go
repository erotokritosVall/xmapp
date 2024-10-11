package http

import (
	"github.com/erotokritosVall/xmapp/pkg/mongo"
	"github.com/erotokritosVall/xmapp/pkg/redis"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type configuration struct {
	AppHost        string   `envconfig:"APP_HOST"`
	AppPort        string   `envconfig:"APP_PORT"`
	AppEnv         string   `envconfig:"APP_ENV"`
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS"`
	RedisConfig    *redis.Configuration
	MongoConfig    *mongo.Configuration
}

func (s *server) readConfiguration() {
	cfg := &configuration{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatal().Msgf("failed to read configuration: %+v", err)
	}

	s.config = cfg
}
