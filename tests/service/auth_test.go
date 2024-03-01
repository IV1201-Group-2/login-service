package service_test

import (
	"strconv"
	"testing"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/model"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/stretchr/testify/require"
)

// Tests that authenticating as a valid user works.
func TestAuthenticateUser(t *testing.T) {
	t.Parallel()

	repository := database.NewUserRepository(tests.Database)

	// Authenticate as applicant
	user, err := service.AuthenticateUser(repository, tests.MockApplicant.Email, tests.MockPassword, &tests.MockApplicant.Role)
	require.NoError(t, err)
	require.Equal(t, user.ID, tests.MockApplicant.ID)
	require.Equal(t, user.Email, tests.MockApplicant.Email)
	require.Equal(t, user.Username, tests.MockApplicant.Username)
	require.Equal(t, user.Role, tests.MockApplicant.Role)

	// Authenticate as recruiter
	user, err = service.AuthenticateUser(repository, tests.MockRecruiter.Username, tests.MockPassword, &tests.MockRecruiter.Role)
	require.NoError(t, err)
	require.Equal(t, user.ID, tests.MockRecruiter.ID)
	require.Equal(t, user.Email, tests.MockRecruiter.Email)
	require.Equal(t, user.Username, tests.MockRecruiter.Username)
	require.Equal(t, user.Role, tests.MockRecruiter.Role)
}

// Tests that authenticating with an invalid identity or role doesn't work.
func TestAuthenticateWrongIdentity(t *testing.T) {
	t.Parallel()

	repository := database.NewUserRepository(tests.Database)

	// Authenticate using the wrong identity
	_, err := service.AuthenticateUser(repository, "wrong", tests.MockPassword, &tests.MockApplicant.Role)
	require.ErrorIs(t, err, service.ErrWrongIdentity)

	// Authenticate using the user's ID
	_, err = service.AuthenticateUser(repository, strconv.Itoa(tests.MockApplicant.ID), tests.MockPassword, &tests.MockApplicant.Role)
	require.ErrorIs(t, err, service.ErrWrongIdentity)

	// Authenticate using the wrong role
	_, err = service.AuthenticateUser(repository, tests.MockApplicant.Email, tests.MockPassword, &tests.MockRecruiter.Role)
	require.ErrorIs(t, err, service.ErrWrongIdentity)
}

// Tests that authenticating with a hashed or invalid password doesn't work.
func TestAuthenticateWrongPassword(t *testing.T) {
	t.Parallel()

	repository := database.NewUserRepository(tests.Database)

	// Authenticate using the wrong password
	_, err := service.AuthenticateUser(repository, tests.MockApplicant.Email, "wrong", &tests.MockApplicant.Role)
	require.ErrorIs(t, err, service.ErrWrongPassword)
	// Authenticate using the user's hashed password
	_, err = service.AuthenticateUser(repository, tests.MockApplicant.Email, tests.MockApplicant.Password, &tests.MockApplicant.Role)
	require.ErrorIs(t, err, service.ErrWrongPassword)
}

// Tests that resetting the password of a user works and the user is modified in the repository.
func TestResetPassword(t *testing.T) {
	t.Parallel()

	repository := database.NewUserRepository(tests.Database)

	// Generate a new random password every time the test is run.
	newPassword := tests.RandomStr(16)
	// Create a new reset token for the user
	claims := model.UserClaims{
		CustomClaims: model.CustomClaims{
			Usage: model.TokenUsageReset,
		},
		User: tests.MockApplicant4,
	}

	// Try before password reset
	_, err := service.AuthenticateUser(repository, tests.MockApplicant4.Email, newPassword, &tests.MockApplicant4.Role)
	require.ErrorIs(t, err, service.ErrMissingPassword)

	err = service.UpdatePassword(repository, claims, newPassword)
	require.NoError(t, err)

	// Try again after password reset
	_, err = service.AuthenticateUser(repository, tests.MockApplicant4.Email, newPassword, &tests.MockApplicant4.Role)
	require.NoError(t, err)
}

// Tests that resetting the password with login claims doesn't work.
func TestResetWrongUsage(t *testing.T) {
	t.Parallel()

	repository := database.NewUserRepository(tests.Database)

	// Generate a new random password every time the test is run.
	newPassword := tests.RandomStr(16)
	// Create a new login token for the user
	claims := model.UserClaims{
		CustomClaims: model.CustomClaims{
			Usage: model.TokenUsageLogin,
		},
		User: tests.MockApplicant4,
	}

	err := service.UpdatePassword(repository, claims, newPassword)
	require.ErrorIs(t, err, service.ErrWrongUsage)
}
