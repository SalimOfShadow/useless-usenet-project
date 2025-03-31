package store

import (
	"context"
	"database/sql"
)

type Subscription struct {
	UserID    int64  `json:"user_id"`
	Newsgroup string `json:"newsgroup"`
}

type SubscriptionStore struct {
	db *sql.DB
}

func (s *SubscriptionStore) Subscribe(ctx context.Context, userID int64, newsgroup string) error {
	query := `INSERT INTO subscriptions (user_id, newsgroup) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, newsgroup)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionStore) Unsubscribe(ctx context.Context, userID int64, newsgroup string) error {
	query := `DELETE FROM subscriptions WHERE user_id = $1 AND newsgroup = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, userID, newsgroup)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *SubscriptionStore) GetSubscriptions(ctx context.Context, userID int64) ([]Newsgroup, error) {
	query := `
		SELECT n.name, n.description, n.created_at
		FROM subscriptions s
		JOIN newsgroups n ON s.newsgroup = n.name
		WHERE s.user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
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