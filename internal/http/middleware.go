package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/erotokritosVall/xmapp/pkg/api"
	"github.com/erotokritosVall/xmapp/pkg/util"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "
)

func (server *server) enableMiddleware() {
	server.router.Use(middleware.Timeout(api.Timeout))
	server.router.Use(middleware.StripSlashes)
	server.router.Use(middleware.GetHead)
	server.router.Use(middleware.RequestID)
	server.router.Use(httplog.RequestLogger(log.Logger))
	server.router.Use(middleware.Recoverer)
	server.router.Use(render.SetContentType(render.ContentTypeJSON))
	server.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   server.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

func (server *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authHeaderKey)
		token := strings.TrimPrefix(bearerPrefix, authHeader)

		if util.IsEmptyOrWhitespace(token) {
			api.Response(http.StatusUnauthorized).Send(w, r)
			return
		}

		loggedOutToken, err := server.redis.GetString(r.Context(), token)
		if err != nil {
			log.Err(err).Msg("could not read redis")
			api.Response(http.StatusUnauthorized).Send(w, r)
			return
		}

		if loggedOutToken != nil && !util.IsEmptyOrWhitespace(*loggedOutToken) {
			api.Response(http.StatusUnauthorized).Send(w, r)
			return
		}

		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return "token", nil
		})

		if !parsedToken.Valid {
			api.Response(http.StatusUnauthorized).Send(w, r)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			log.Err(err).Msg("could not parse claims")
			api.Response(http.StatusUnauthorized).Send(w, r)
			return
		}

		expTime, err := claims.GetExpirationTime()
		if err != nil {
			log.Err(err).Msg("could not read claims expiration time")
			api.Response(http.StatusUnauthorized).Send(w, r)
			return
		}

		if expTime.Before(time.Now().UTC()) {
			api.Response(http.StatusUnauthorized).Send(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "auth_token", token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
