package database_test

import (
	"strconv"
	"testing"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/stretchr/testify/require"
)

// Test that a user can be queried from the database.
func TestQueryUser(t *testing.T) {
	t.Parallel()

	repository := database.NewUserRepository(tests.Database)

	// Query for applicant
	applicant, err := repository.Query(tests.MockApplicant.Email)
	require.NoError(t, err)

	require.Equal(t, applicant.ID, tests.MockApplicant.ID)
	require.Equal(t, applicant.Username, tests.MockApplicant.Username)
	require.Equal(t, applicant.Email, tests.MockApplicant.Email)
	require.Equal(t, applicant.Role, tests.MockApplicant.Role)

	require.Equal(t, applicant.Password, tests.MockPasswordBcrypt)
	require.NotEqual(t, applicant.Password, tests.MockPassword)

	// Query for recruiter
	recruiter, err := repository.Query(tests.MockRecruiter.Username)
	require.NoError(t, err)

	require.Equal(t, recruiter.ID, tests.MockRecruiter.ID)
	require.Equal(t, recruiter.Username, tests.MockRecruiter.Username)
	require.Equal(t, recruiter.Email, tests.MockRecruiter.Email)
	require.Equal(t, recruiter.Role, tests.MockRecruiter.Role)

	require.Equal(t, recruiter.Password, tests.MockPasswordBcrypt)
	require.NotEqual(t, recruiter.Password, tests.MockPassword)

	// Query for empty identity
	// NOTE: This should NOT fail with our test data. The service layer guards against information leaks.
	// It's not possible to handle in the database layer because there are some users
	// thave have empty username or emails.
	user, err := repository.Query("")
	require.NotNil(t, user)
	require.NoError(t, err)
}

// Test that missing users return "no user found" from database.
func TestQueryInvalidIdentity(t *testing.T) {
	t.Parallel()

	repository := database.NewUserRepository(tests.Database)

	// Query for a user ID
	user, err := repository.Query(strconv.Itoa(tests.MockApplicant.ID))
	require.Nil(t, user)
	require.ErrorIs(t, err, database.ErrUserNotFound)

	// Query for invalid identity
	user, err = repository.Query("wrong")
	require.Nil(t, user)
	require.ErrorIs(t, err, database.ErrUserNotFound)
}

// Test that the password of a user can be changed.
func TestResetPassword(t *testing.T) {
	t.Parallel()

	// Generate a new random password every time the test is run.
	newPassword := tests.RandomStr(16)
	repository := database.NewUserRepository(tests.Database)

	// Query for the user once
	user, err := repository.Query(tests.MockApplicant5.Email)
	require.NoError(t, err)
	require.Empty(t, user.Password)

	err = repository.UpdatePassword(tests.MockApplicant5.ID, newPassword)
	require.NoError(t, err)

	// Query for the user again
	user, err = repository.Query(tests.MockApplicant5.Email)
	require.NoError(t, err)
	require.Equal(t, user.Password, newPassword)
}
