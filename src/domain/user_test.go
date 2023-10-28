package domain_test

import (
	"testing"

	"github.com/apm-dev/oha/src/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_User_IsNew(t *testing.T) {
	testCases := []struct {
		desc   string
		user   *domain.User
		expect bool
	}{
		{
			desc: "when user has ID, should return false",
			user: &domain.User{
				ID:   uuid.New().String(),
				Name: "Amir",
			},
			expect: false,
		},
		{
			desc: "when user does not have ID, should return true",
			user: &domain.User{
				ID:   "",
				Name: "Amir",
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expect, tc.user.IsNew(), tc.desc)
	}
}

func Test_User_TryGenerateID(t *testing.T) {
	testCases := []struct {
		desc   string
		user   *domain.User
		expect bool
	}{
		{
			desc: "when user has ID, should return false",
			user: &domain.User{
				ID:   uuid.New().String(),
				Name: "Amir",
			},
			expect: false,
		},
		{
			desc: "when user does not have ID, should fill the ID and return true",
			user: &domain.User{
				ID:   "",
				Name: "Amir",
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expect, tc.user.TryGenerateID(), tc.desc)
		if tc.expect == true {
			assert.NotEmpty(t, tc.user.ID, tc.desc)
		}
	}
}
