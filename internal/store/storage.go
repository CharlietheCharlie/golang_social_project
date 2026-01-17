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
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error

		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error // practice transaction
		Activate(ctx context.Context, token string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
