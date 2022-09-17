package repository_mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-seidon/local/internal/datetime"
	"github.com/go-seidon/local/internal/repository"
)

type RepositoryParam struct {
	mClient *sql.DB
	rClient *sql.DB
	clock   datetime.Clock
}

type RepoOption = func(*RepositoryParam)

func WithDbMaster(dbClient *sql.DB) RepoOption {
	return func(ro *RepositoryParam) {
		ro.mClient = dbClient
	}
}

func WithDbReplica(dbClient *sql.DB) RepoOption {
	return func(ro *RepositoryParam) {
		ro.rClient = dbClient
	}
}

func WithClock(clock datetime.Clock) RepoOption {
	return func(ro *RepositoryParam) {
		ro.clock = clock
	}
}

type provider struct {
	authRepo *authRepository
	fileRepo *fileRepository
}

func (p *provider) Init(ctx context.Context) error {
	return nil
}

func (p *provider) GetAuthRepo() repository.AuthRepository {
	return p.authRepo
}

func (p *provider) GetFileRepo() repository.FileRepository {
	return p.fileRepo
}

func NewRepository(opts ...RepoOption) (*provider, error) {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if p.mClient == nil {
		return nil, fmt.Errorf("invalid db client specified")
	}
	if p.rClient == nil {
		return nil, fmt.Errorf("invalid db client specified")
	}

	var clock datetime.Clock
	if p.clock == nil {
		clock = datetime.NewClock()
	} else {
		clock = p.clock
	}

	authRepo := &authRepository{
		mClient: p.mClient,
		rClient: p.rClient,
		clock:   clock,
	}
	fileRepo := &fileRepository{
		mClient: p.mClient,
		rClient: p.rClient,
		clock:   clock,
	}

	repo := &provider{
		authRepo: authRepo,
		fileRepo: fileRepo,
	}
	return repo, nil
}
