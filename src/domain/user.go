package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserService interface {
	AddNewUser(ctx context.Context, name string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type UserRepository interface {
	Insert(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
}

// Domain Logics
func NewUser(name string) (*User, error) {
	if len(name) < 3 {
		return nil, errors.Wrap(ErrInvalidArgument, "'name' length should be at least 3 characters")
	}
	user := &User{
		Name: name,
	}
	if ok := user.TryGenerateID(); !ok {
		return nil, errors.Wrap(ErrInternalServer, "failed to generate ID for user")
	}
	return user, nil
}

func (u *User) IsNew() bool {
	return len(u.ID) <= 0
}

func (u *User) TryGenerateID() bool {
	if u.IsNew() {
		u.ID = uuid.New().String()
		return true
	}
	return false
}
