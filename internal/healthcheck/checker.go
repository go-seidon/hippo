package healthcheck

import (
	"context"
	"time"
)

type DataSource interface {
	Ping(ctx context.Context) error
}

type repoPingChecker struct {
	dataSource DataSource
}

type RepoPingResult struct {
	CheckedAt time.Time
}

func (p *repoPingChecker) Status() (interface{}, error) {
	ctx := context.Background()
	err := p.dataSource.Ping(ctx)
	if err != nil {
		return nil, err
	}

	ts := time.Now()
	res := &RepoPingResult{ts}
	return res, nil
}

func NewRepoPingChecker(dataSource DataSource) *repoPingChecker {
	return &repoPingChecker{dataSource}
}
