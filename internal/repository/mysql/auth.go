package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/datetime"
	db_mysql "github.com/go-seidon/provider/mysql"
)

type authRepository struct {
	mClient db_mysql.Client
	rClient db_mysql.Client
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
		return nil, repository.ErrNotFound
	}
	return nil, err
}

func NewAuth(opts ...RepoOption) *authRepository {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	clock := p.clock
	if clock == nil {
		clock = datetime.NewClock()
	}

	return &authRepository{
		mClient: p.mClient,
		rClient: p.rClient,
		clock:   clock,
	}
}
