package repository_mysql

import (
	"fmt"

	db_mysql "github.com/go-seidon/hippo/internal/db-mysql"
	"github.com/go-seidon/provider/datetime"
)

type RepositoryParam struct {
	mClient db_mysql.Client
	rClient db_mysql.Client
	clock   datetime.Clock
}

type RepoOption = func(*RepositoryParam)

func WithDbMaster(dbClient db_mysql.Client) RepoOption {
	return func(ro *RepositoryParam) {
		ro.mClient = dbClient
	}
}

func WithDbReplica(dbClient db_mysql.Client) RepoOption {
	return func(ro *RepositoryParam) {
		ro.rClient = dbClient
	}
}

func WithClock(clock datetime.Clock) RepoOption {
	return func(ro *RepositoryParam) {
		ro.clock = clock
	}
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

	clock := p.clock
	if clock == nil {
		clock = datetime.NewClock()
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
		mClient:  p.mClient,
		rClient:  p.rClient,
		authRepo: authRepo,
		fileRepo: fileRepo,
	}
	return repo, nil
}
