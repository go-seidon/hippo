package repository_mysql

import (
	"database/sql"

	"github.com/go-seidon/local/internal/datetime"
)

type RepositoryOption struct {
	mClient *sql.DB
	rClient *sql.DB
	clock   datetime.Clock
}

type RepoOption = func(*RepositoryOption)

func WithDbMaster(dbClient *sql.DB) RepoOption {
	return func(ro *RepositoryOption) {
		ro.mClient = dbClient
	}
}

func WithDbReplica(dbClient *sql.DB) RepoOption {
	return func(ro *RepositoryOption) {
		ro.rClient = dbClient
	}
}

func WithClock(clock datetime.Clock) RepoOption {
	return func(ro *RepositoryOption) {
		ro.clock = clock
	}
}
