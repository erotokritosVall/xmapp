package application

import (
	"context"
	"time"

	domain "github.com/erotokritosVall/xmapp/internal/users/domain"
	"github.com/erotokritosVall/xmapp/pkg/errors"
	"github.com/erotokritosVall/xmapp/pkg/redis"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 14
	tokenTtl = 30 * time.Minute
)

type userService struct {
	repo  domain.UserRepository
	redis redis.Redis
}

func NewService(repo domain.UserRepository,
	redis redis.Redis) domain.UserService {
	return &userService{
		repo:  repo,
		redis: redis,
	}
}

func (srv *userService) Create(ctx context.Context, user *domain.User) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), hashCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)

	return srv.repo.Create(ctx, user)
}

func (srv *userService) Login(ctx context.Context, email, password string) (*string, error) {
	user, err := srv.repo.ReadByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.ErrInvalidPassword
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "api",
		"u":   user.Id.Hex(),
		"exp": time.Now().UTC().Add(tokenTtl),
	})

	token, err := t.SignedString("token")
	return &token, err
}

func (srv *userService) Logout(ctx context.Context) error {
	token := ctx.Value("auth_token").(string)

	return srv.redis.SetString(ctx, token, "1", tokenTtl)
}
