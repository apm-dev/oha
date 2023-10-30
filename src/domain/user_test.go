package domain_test

import (
	"testing"

	"github.com/apm-dev/oha/src/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewUser(t *testing.T) {
	testCases := []struct {
		desc          string
		inputName     string
		expectedUser  *domain.User
		expectedError error
	}{
		{
			desc:          "Valid user",
			inputName:     "Amir",
			expectedUser:  &domain.User{Name: "Amir"},
			expectedError: nil,
		},
		{
			desc:          "Invalid name",
			inputName:     "A",
			expectedUser:  nil,
			expectedError: domain.ErrInvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := domain.NewUser(tc.inputName)

			if tc.expectedUser != nil {
				assert.NotEmpty(t, user.ID, tc.desc)
				assert.Equal(t, tc.expectedUser.Name, user.Name, tc.desc)
			}
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError, tc.desc)
			} else {
				assert.NoError(t, err, tc.desc)
			}
		})
	}
}

func Test_IsNew(t *testing.T) {
	user := &domain.User{}
	assert.True(t, user.IsNew())

	user.ID = "123-abc-456-def"
	assert.False(t, user.IsNew())
}

func Test_TryGenerateID(t *testing.T) {
	user := &domain.User{}
	generated := user.TryGenerateID()

	assert.True(t, generated)
	assert.NotEmpty(t, user.ID)
	_, err := uuid.Parse(user.ID)
	assert.NoError(t, err)

	// ID already exists, should not generate a new one
	previousID := user.ID
	generated = user.TryGenerateID()
	assert.False(t, generated)
	assert.Equal(t, previousID, user.ID)
}
