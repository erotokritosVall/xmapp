package infrastructure

import (
	"context"

	domain "github.com/erotokritosVall/xmapp/internal/users/domain"
	"github.com/erotokritosVall/xmapp/pkg/errors"
	"github.com/erotokritosVall/xmapp/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "users"

type userRepository struct {
	db *mongo.Database
}

func New(db *mongo.Database) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Create(ctx context.Context, user *domain.User) (*string, error) {
	result, err := repo.db.Collection(collectionName).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID)

	return util.New(id.Hex()), nil
}

func (repo *userRepository) ReadByEmail(ctx context.Context, email string) (*domain.User, error) {
	where := bson.M{
		"email": email,
	}

	user := &domain.User{}

	if err := repo.db.Collection(collectionName).FindOne(ctx, where).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}

		return nil, err
	}

	return user, nil
}
