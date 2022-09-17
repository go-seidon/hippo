package repository_mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-seidon/local/internal/datetime"
	"github.com/go-seidon/local/internal/repository"
)

type authRepository struct {
	mClient *sql.DB
	rClient *sql.DB
	clock   datetime.Clock
}

func (r *authRepository) FindClient(ctx context.Context, p repository.FindClientParam) (*repository.FindClientResult, error) {
	sqlQuery := `
		SELECT 
			client_id, client_secret
		FROM auth_client
		WHERE client_id = ?
	`

	var res repository.FindClientResult
	row := r.rClient.QueryRow(sqlQuery, p.ClientId)
	err := row.Scan(
		&res.ClientId,
		&res.ClientSecret,
	)
	if err == nil {
		return &res, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrorRecordNotFound
	}
	return nil, err
}

func NewAuthRepository(opts ...RepoOption) (*authRepository, error) {
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

	r := &authRepository{
		mClient: p.mClient,
		rClient: p.rClient,
		clock:   clock,
	}
	return r, nil
}
