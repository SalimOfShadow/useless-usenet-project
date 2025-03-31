package store

import (
	"context"
	"database/sql"
	"errors"
)

type Newsgroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type NewsgroupStore struct {
	db *sql.DB
}

func (s *NewsgroupStore) List(ctx context.Context) ([]Newsgroup, error) {
	query := `SELECT name, description, created_at FROM newsgroups ORDER BY name ASC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsgroups []Newsgroup
	for rows.Next() {
		var n Newsgroup
		err := rows.Scan(&n.Name, &n.Description, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		newsgroups = append(newsgroups, n)
	}

	return newsgroups, nil
}

func (s *NewsgroupStore) GetByName(ctx context.Context, name string) (*Newsgroup, error) {
	query := `SELECT name, description, created_at FROM newsgroups WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var newsgroup Newsgroup
	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&newsgroup.Name,
		&newsgroup.Description,
		&newsgroup.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &newsgroup, nil
}

func (s *NewsgroupStore) Create(ctx context.Context, newsgroup *Newsgroup) error {
	query := `INSERT INTO newsgroups (name, description) VALUES ($1, $2) RETURNING created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, newsgroup.Name, newsgroup.Description).Scan(&newsgroup.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *NewsgroupStore) Delete(ctx context.Context, name string) error {
	query := `DELETE FROM newsgroups WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, name)
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
