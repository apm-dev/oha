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
		input        *Inputs
		prepareMocks func(m *Mocks)
		expectedUser *domain.User
		expectedErr  error
	}{
		{
			desc: "Invalid name with length < 3",
			input: &Inputs{
				name: []string{"", "a", "ab"}[rand.Intn(3)],
			},
			prepareMocks: func(m *Mocks) {},
			expectedUser: nil,
			expectedErr:  domain.ErrInvalidArgument,
		},
		{
			desc: "Internal server error",
			input: &Inputs{
				name: "Amir",
			},
			prepareMocks: func(m *Mocks) {
				m.ur.On("Insert", mock.Anything, mock.Anything).
					Return(errors.New("some db error")).Once()
			},
			expectedUser: nil,
			expectedErr:  domain.ErrInternalServer,
		},
		{
			desc: "Create user",
			input: &Inputs{
				name: "Parsa",
			},
			prepareMocks: func(m *Mocks) {
				m.ur.On("Insert", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			expectedUser: &domain.User{
				Name: "Parsa",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Arrange
			mocks := &Mocks{
				ur: mocks.NewUserRepository(t),
			}
			tc.prepareMocks(mocks)
			svc := user.NewService(mocks.ur)
			// Action
			user, err := svc.AddNewUser(context.Background(), tc.input.name)
			// Assert
			if tc.expectedUser != nil {
				assert.NotEmpty(t, user.ID, tc.desc)
				assert.Equal(t, tc.expectedUser.Name, user.Name, tc.desc)
			}
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr, tc.desc)
			} else {
				assert.NoError(t, err, tc.desc)
			}
			mocks.ur.AssertExpectations(t)
		})
	}
}

func Test_GetUserByID(t *testing.T) {
	testCases := []struct {
		desc         string
		userID       string
		repoResponse *domain.User
		repoErr      error
		expectedUser *domain.User
		expectedErr  error
	}{
		{
			desc:         "User found",
			userID:       "abc",
			repoResponse: &domain.User{ID: "abc", Name: "Amir"},
			repoErr:      nil,
			expectedUser: &domain.User{ID: "abc", Name: "Amir"},
			expectedErr:  nil,
		},
		{
			desc:         "User not found",
			userID:       "def",
			repoResponse: nil,
			repoErr:      nil,
			expectedUser: nil,
			expectedErr:  domain.ErrNotFound,
		},
		{
			desc:         "Internal server error",
			userID:       "ghi",
			repoResponse: nil,
			repoErr:      errors.New("some error"),
			expectedUser: nil,
			expectedErr:  domain.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Arrange
			repo := mocks.NewUserRepository(t)
			repo.On("FindByID", mock.Anything, tc.userID).Return(tc.repoResponse, tc.repoErr)

			s := user.NewService(repo)
			// Action
			user, err := s.GetUserByID(context.Background(), tc.userID)
			// Assert
			assert.Equal(t, tc.expectedUser, user, tc.desc)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr, tc.desc)
			} else {
				assert.NoError(t, err, tc.desc)
			}
			repo.AssertExpectations(t)
		})
	}
}
