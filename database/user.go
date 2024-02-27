package database

import (
	"database/sql"
	"errors"

	"github.com/IV1201-Group-2/login-service/model"
	sq "github.com/Masterminds/squirrel"
)

// Query the database for a user with the specified identity.
func QueryUser(identity string) (*model.User, error) {
	var name, email, password sql.NullString
	var user model.User

	// Begin transaction:
	// If user is spread across multiple tables all reads need to be done at the same time.
	tx, err := connection.Begin()
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

	if err = tx.Commit(); err != nil {
		return nil, ErrQueryFailed.Wrap(err)
	}

	// Once transaction has been committed, we can read all "potentially null" strings.
	user.Username = name.String
	user.Email = email.String
	user.Password = password.String

	return &user, nil
}

// Update the password for a user with the specified ID.
func UpdatePassword(id int, hashedPassword string) error {
	// Begin transaction:
	// If user is spread across multiple tables all writes need to be done at the same time.
	tx, err := connection.Begin()
	if err != nil {
		return ErrQueryFailed.Wrap(err)
	}
	// Transaction will be automatically rolled back if the function returns an error.
	defer tx.Rollback()

	query := sq.StatementBuilder.RunWith(tx).
		Update("person").
		Set("password", hashedPassword).
		Where(sq.Eq{"person_id": id})

	result, err := query.Exec()
	if err != nil {
		return ErrQueryFailed.Wrap(err)
	}
	if rows, err := result.RowsAffected(); err != nil || rows != 1 {
		return ErrUserNotFound.Wrap(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrQueryFailed.Wrap(err)
	}

	return nil
}
