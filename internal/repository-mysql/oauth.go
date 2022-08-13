package repository_mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-seidon/local/internal/datetime"
	"github.com/go-seidon/local/internal/repository"
)

type oAuthRepository struct {
	mClient *sql.DB
	rClient *sql.DB
	clock   datetime.Clock
}

// @note: use master client to avoid
// unable to modify state when replica database is down
func (r *oAuthRepository) FindClient(ctx context.Context, p repository.FindClientParam) (*repository.FindClientResult, error) {
	sqlQuery := `
		SELECT 
			client_id, client_secret
		FROM oauth_client
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

func NewOAuthRepository(opts ...RepoOption) (*oAuthRepository, error) {
	option := RepositoryOption{}
	for _, opt := range opts {
		opt(&option)
	}

	if option.mClient == nil {
		return nil, fmt.Errorf("invalid db client specified")
	}
	if option.rClient == nil {
		return nil, fmt.Errorf("invalid db client specified")
	}

	var clock datetime.Clock
	if option.clock == nil {
		clock = datetime.NewClock()
	} else {
		clock = option.clock
	}

	r := &oAuthRepository{
		mClient: option.mClient,
		rClient: option.rClient,
		clock:   clock,
	}
	return r, nil
}
