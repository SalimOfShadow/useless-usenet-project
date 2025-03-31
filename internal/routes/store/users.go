package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Username string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, email, username, created_at FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, username, created_at FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Authenticate(ctx context.Context, email, password string) (*User, error) {
	query := `SELECT id, email, username, created_at FROM users WHERE email = $1 AND password = crypt($2, password)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, email, password).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) Create(ctx context.Context, transaction *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (email, password, username, created_at)
		VALUES ($1, crypt($2, gen_salt('bf')), $3, NOW())
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var err error
	if transaction != nil {
		err = transaction.QueryRowContext(ctx, query, user.Email, user.Password, user.Username).Scan(&user.ID, &user.CreatedAt)
	} else {
		err = s.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Username).Scan(&user.ID, &user.CreatedAt)
	}

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint" { // Adjust based on actual error message
			return ErrConflict
		}
		return err
	}

	return nil
}