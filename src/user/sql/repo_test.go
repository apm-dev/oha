package usersql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/apm-dev/oha/src/domain"
	"github.com/apm-dev/oha/src/user/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserRepo_Insert(t *testing.T) {
	// Create a new SQL mock database connection and repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := usersql.NewUserRepo(db)

	// Test case 1: Successful insert
	user := &domain.User{ID: "d06f3b0f-16a5-4414-a05e-92e09d170cc9", Name: "Amir"}
	mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
		WithArgs(user.ID, user.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(context.Background(), user)
	require.NoError(t, err)

	// Test case 2: Insert error
	mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
		WithArgs(user.ID, user.Name).
		WillReturnError(sql.ErrNoRows)

	err = repo.Insert(context.Background(), user)
	require.Error(t, err)

	// Ensure expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_UserRepo_FindByID(t *testing.T) {
	// Create a new SQL mock database connection and repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := usersql.NewUserRepo(db)

	// Test case 1: User found
	user := &domain.User{ID: "d06f3b0f-16a5-4414-a05e-92e09d170cc9", Name: "Parsa"}
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(user.ID, user.Name, user.CreatedAt, user.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\$1").
		WithArgs(user.ID).
		WillReturnRows(rows)

	result, err := repo.FindByID(context.Background(), user.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user, result)

	// Test case 2: User not found
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=\\$1").
		WithArgs("not-existing-id").
		WillReturnError(sql.ErrNoRows)

	result, err = repo.FindByID(context.Background(), "not-existing-id")
	require.NoError(t, err)
	assert.Nil(t, result)

	// Ensure expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
