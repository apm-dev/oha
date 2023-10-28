package user

import (
	"context"

	"github.com/apm-dev/oha/src/domain"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	userRepo domain.UserRepository
}

func NewService(
	ur domain.UserRepository,
) domain.UserService {
	return &Service{
		userRepo: ur,
	}
}

func (s *Service) AddNewUser(ctx context.Context, name string) (*domain.User, error) {
	const op = "user.service.AddNewUser"

	user, err := domain.NewUser(name)
	if !errors.Is(err, domain.ErrInvalidArgument) {
		log.Errorf("%s failed to create a user, err: '%s'", op, err)
	}

	err = s.userRepo.Save(ctx, user)
	if err != nil {
		log.Errorf("%s failed to persist user, err: '%s'", op, err)
		return nil, errors.Wrap(domain.ErrInternalServer, "failed to add the user")
	}
	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	panic("not implemented") // TODO: Implement
}
