package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	usersApp "github.com/erotokritosVall/xmapp/internal/users/application"
	usersInfra "github.com/erotokritosVall/xmapp/internal/users/infrastructure"
	"github.com/erotokritosVall/xmapp/pkg/api"
	mg "github.com/erotokritosVall/xmapp/pkg/mongo"
	"github.com/erotokritosVall/xmapp/pkg/redis"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	apps   []api.App
	redis  redis.Redis
	db     *mongo.Database
	router chi.Router
	config *configuration
}

func New() *server {
	s := &server{}

	s.readConfiguration()
	s.startupRedis()
	s.startupMongo()
	s.enableRouting()
	s.enableMiddleware()
	s.initializeApps()
	s.registerEndpoints()

	return s
}

func (s *server) Start(exitChannel chan os.Signal) {
	go func() {
		addr := fmt.Sprintf("%s:%s", s.config.AppHost, s.config.AppPort)
		if err := http.ListenAndServe(addr, s.router); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()

	log.Debug().Msg("server started")

	<-exitChannel
}

func (s *server) startupRedis() {
	r, err := redis.New(context.Background(), s.config.RedisConfig)
	if err != nil {
		log.Fatal().Msgf("failed to startup redis: %+v", err)
	}

	s.redis = r
}

func (s *server) startupMongo() {
	m, err := mg.New(context.Background(), s.config.MongoConfig)
	if err != nil {
		log.Fatal().Msgf("failed to startup mongo: %+v", err)
	}

	s.db = m
}

func (s *server) enableRouting() {
	s.router = chi.NewRouter()
}

func (s *server) initializeApps() {
	s.apps = []api.App{}

	userRepo := usersInfra.New(s.db)
	userSrv := usersApp.NewService(userRepo, s.redis)
	s.apps = append(s.apps, usersApp.NewApp(userSrv))
}

func (s *server) registerEndpoints() {
	for _, app := range s.apps {
		app.RegisterPublicEndpoints(s.router)
	}

	s.router.Group(func(r chi.Router) {
		r.Use(s.authMiddleware)

		for _, app := range s.apps {
			app.RegisterProtectedEndpoints(r)
		}
	})
}
