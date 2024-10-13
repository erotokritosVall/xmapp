package util

import (
	"github.com/erotokritosVall/xmapp/pkg/mongo"
)

type TestsConfig struct {
	AppHost     string `envconfig:"APP_HOST"`
	AppPort     string `envconfig:"APP_PORT"`
	MongoConfig *mongo.Configuration

	// Should be added in secrets
	TestEmail string `envconfig:"TEST_EMAIL"`
	TestPass  string `envconfig:"TEST_PASS"`
}
