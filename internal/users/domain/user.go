package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*string, error)
	ReadByEmail(ctx context.Context, email string) (*User, error)
}

type UserService interface {
	Create(ctx context.Context, user *User) (*string, error)
	Login(ctx context.Context, email, password string) (*string, error)
	Logout(ctx context.Context) error
}
