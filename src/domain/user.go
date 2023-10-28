package domain

import (
	"context"

	"github.com/google/uuid"
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
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
}

// Domain Logics
func (u *User) IsNew() bool {
	return len(u.ID) <= 0
}

func (u *User) TryGenerateId() bool {
	if u.IsNew() {
		u.ID = uuid.New().String()
		return true
	}
	return false
}
