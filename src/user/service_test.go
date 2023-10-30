package user_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/apm-dev/oha/src/domain"
	"github.com/apm-dev/oha/src/domain/mocks"
	"github.com/apm-dev/oha/src/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AddNewUser(t *testing.T) {
	type Mocks struct {
		ur *mocks.UserRepository
	}
	type Inputs struct {
		name string
	}
	testCases := []struct {
		desc         string
		inputs       *Inputs
		prepareMocks func(m *Mocks)
		expectResult *domain.User
		expectErr    error
	}{
		{
			desc: "when receives invalid name with length < 3, should return error",
			inputs: &Inputs{
				name: []string{"", "a", "ab"}[rand.Intn(3)],
			},
			prepareMocks: func(m *Mocks) {},
			expectResult: nil,
			expectErr:    domain.ErrInvalidArgument,
		},
		{
			desc: "when fails to persist user to DB, should return error",
			inputs: &Inputs{
				name: "Amir",
			},
			prepareMocks: func(m *Mocks) {
				m.ur.On("Insert", mock.Anything, mock.Anything).
					Return(errors.New("some db error")).Once()
			},
			expectResult: nil,
			expectErr:    domain.ErrInternalServer,
		},
		{
			desc: "when creates and persists user to DB, should return user",
			inputs: &Inputs{
				name: "Parsa",
			},
			prepareMocks: func(m *Mocks) {
				m.ur.On("Insert", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			expectResult: &domain.User{
				Name: "Parsa",
			},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		// Arrange
		mocks := &Mocks{
			ur: mocks.NewUserRepository(t),
		}
		tc.prepareMocks(mocks)
		svc := user.NewService(mocks.ur)

		// Action
		user, err := svc.AddNewUser(context.Background(), tc.inputs.name)

		// Assert
		if tc.expectErr != nil {
			assert.ErrorIs(t, err, tc.expectErr, tc.desc)
		} else {
			assert.NoError(t, err, tc.desc)
		}
		if tc.expectResult != nil {
			assert.NotEmpty(t, user.ID, tc.desc)
			assert.Equal(t, tc.expectResult.Name, user.Name, tc.desc)
		}
		mocks.ur.AssertExpectations(t)
	}
}
