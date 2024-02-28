package database

import (
	"database/sql"
	"errors"

	"github.com/IV1201-Group-2/login-service/model"
	sq "github.com/Masterminds/squirrel"
)

type UserRepository struct {
	conn *sql.DB
}

// NewUserRepository creates a new repository from a database connection.
func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn}
}

// Query the repository for a user with the specified identity.
func (u *UserRepository) Query(identity string) (*model.User, error) {
	var name, email, password sql.NullString
	var user model.User

	// Begin transaction:
	// If user is spread across multiple tables all reads need to be done at the same time.
	tx, err := u.conn.Begin()
	if err != nil {
		return nil, ErrQueryFailed.Wrap(err)
	}
	// Transaction will be automatically rolled back if the function returns an error.
	defer tx.Rollback()

	query := sq.StatementBuilder.RunWith(tx).
		Select("person_id", "username", "email", "password", "role_id").
		From("person").
		Where(sq.Or{sq.Eq{"username": identity}, sq.Eq{"email": identity}})

	err = query.Scan(&user.ID, &name, &email, &password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound.Wrap(err)
	} else if err != nil {
		return nil, ErrQueryFailed.Wrap(err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, ErrQueryFailed.Wrap(err)
	}

	// Once transaction has been committed, we can read all "potentially null" strings.
	user.Username = name.String
	user.Email = email.String
	user.Password = password.String

	return &user, nil
}

// Update the password for a user in the repository with the specified ID.
func (u *UserRepository) UpdatePassword(id int, password string) error {
	// Begin transaction:
	// If user is spread across multiple tables all writes need to be done at the same time.
	tx, err := u.conn.Begin()
	if err != nil {
		return ErrQueryFailed.Wrap(err)
	}
	// Transaction will be automatically rolled back if the function returns an error.
	defer tx.Rollback()

	query := sq.StatementBuilder.RunWith(tx).
		Update("person").
		Set("password", password).
		Where(sq.Eq{"person_id": id})

	result, err := query.Exec()
	if err != nil {
		return ErrQueryFailed.Wrap(err)
	}
	// If no rows were affected, the user was not found
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return ErrUserNotFound.Wrap(err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return ErrQueryFailed.Wrap(err)
	}

	return nil
}
