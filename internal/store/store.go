package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)


type Storage struct {
	Newsgroups interface {
		GetByName(ctx context.Context, name string) (*Newsgroup, error)
		List(ctx context.Context) ([]Newsgroup, error)
		Create(ctx context.Context, newsgroup *Newsgroup) error
		Delete(ctx context.Context, name string) error
	}
	Articles interface {
		GetByID(ctx context.Context, id int64) (*Article, error)
		GetByNewsgroup(ctx context.Context, newsgroup string, limit int) ([]ArticleWithMetadata, error) 
		Create(ctx context.Context, article *Article) error
		Delete(ctx context.Context, id int64) error
	}
	Users interface {
		GetByID(ctx context.Context, id int64) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
		Create(ctx context.Context, transaction *sql.Tx, user *User) error
		Authenticate(ctx context.Context, email, password string) (*User, error)
		Delete(ctx context.Context, id int64) error
	}
	Subscriptions interface {
		Subscribe(ctx context.Context, userID int64, newsgroup string) error
		Unsubscribe(ctx context.Context, userID int64, newsgroup string) error
		GetSubscriptions(ctx context.Context, userID int64) ([]Newsgroup, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Newsgroups:   &NewsgroupStore{db},  
		Articles:     &ArticleStore{db},   
		Users:        &UserStore{db},      
		Subscriptions: &SubscriptionStore{db}, 
	}
}
