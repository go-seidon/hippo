package mysql

import (
	"github.com/go-seidon/provider/mysql"
	"gorm.io/gorm"
)

type RepositoryParam struct {
	gormClient *gorm.DB
	dbClient   mysql.Client
}

type RepoOption = func(*RepositoryParam)

func WithGormClient(g *gorm.DB) RepoOption {
	return func(p *RepositoryParam) {
		p.gormClient = g
	}
}

func WithDbClient(c mysql.Client) RepoOption {
	return func(p *RepositoryParam) {
		p.dbClient = c
	}
}
