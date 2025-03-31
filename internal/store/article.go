package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Article struct {
	ID        int64    `json:"id"`
	Subject   string   `json:"subject"`
	Author    string   `json:"author"`
	Newsgroup string   `json:"newsgroup"`
	Body      string   `json:"body"`
	CreatedAt string   `json:"created_at"`
	Tags      []string `json:"tags"`
}

type ArticleWithMetadata struct {
	Article
	ReplyCount int `json:"reply_count"`
}

type ArticleStore struct {
	db *sql.DB
}

func (s *ArticleStore) GetByNewsgroup(ctx context.Context, newsgroup string, limit int) ([]ArticleWithMetadata, error) {
	query := `
		SELECT a.id, a.subject, a.author, a.newsgroup, a.body, a.created_at, a.tags,
		COUNT(r.id) AS reply_count
		FROM articles a
		LEFT JOIN replies r ON r.article_id = a.id
		WHERE a.newsgroup = $1
		GROUP BY a.id
		ORDER BY a.created_at DESC
		LIMIT $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, newsgroup, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []ArticleWithMetadata
	for rows.Next() {
		var a ArticleWithMetadata
		err := rows.Scan(
			&a.ID,
			&a.Subject,
			&a.Author,
			&a.Newsgroup,
			&a.Body,
			&a.CreatedAt,
			pq.Array(&a.Tags),
			&a.ReplyCount,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (s *ArticleStore) Create(ctx context.Context, article *Article) error {
	query := `
		INSERT INTO articles (subject, author, newsgroup, body, tags)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		article.Subject,
		article.Author,
		article.Newsgroup,
		article.Body,
		pq.Array(article.Tags),
	).Scan(
		&article.ID,
		&article.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArticleStore) GetByID(ctx context.Context, id int64) (*Article, error) {
	query := `
		SELECT id, subject, author, newsgroup, body, created_at, tags
		FROM articles
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var article Article
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID,
		&article.Subject,
		&article.Author,
		&article.Newsgroup,
		&article.Body,
		&article.CreatedAt,
		pq.Array(&article.Tags),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &article, nil
}

func (s *ArticleStore) Delete(ctx context.Context, articleID int64) error {
	query := `DELETE FROM articles WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, articleID)
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
