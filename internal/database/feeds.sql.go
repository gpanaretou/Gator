// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, name, url, user_id, created_at, updated_at, last_fetched_at
`

type CreateFeedParams struct {
	ID        uuid.UUID
	Name      string
	Url       string
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeed = `-- name: GetFeed :one
SELECT id, name, url, user_id, created_at, updated_at, last_fetched_at FROM feeds
WHERE url = $1 LIMIT 1
`

func (q *Queries) GetFeed(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT id, name, url, user_id, created_at, updated_at, last_fetched_at FROM feeds
ORDER BY user_id
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT id, name, url, user_id, created_at, updated_at, last_fetched_at FROM feeds
ORDER BY last_fetched_at NULLS FIRST
`

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const getTotalNumberOfFeeds = `-- name: GetTotalNumberOfFeeds :one
SELECT COUNT(*) FROM feeds
`

func (q *Queries) GetTotalNumberOfFeeds(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTotalNumberOfFeeds)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const markFeedFetched = `-- name: MarkFeedFetched :one
UPDATE feeds
SET 
    last_fetched_at = $2,
    updated_at = $3
WHERE id = $1
RETURNING id, name, url, user_id, created_at, updated_at, last_fetched_at
`

type MarkFeedFetchedParams struct {
	ID            uuid.UUID
	LastFetchedAt sql.NullTime
	UpdatedAt     time.Time
}

func (q *Queries) MarkFeedFetched(ctx context.Context, arg MarkFeedFetchedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, markFeedFetched, arg.ID, arg.LastFetchedAt, arg.UpdatedAt)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}
