package mysql_test

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func OpenDb(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = "admin:123456@tcp(127.0.0.1:3307)/hippo_test?parseTime=true"
	}
	client, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func RunDbMigration(db *sql.DB) error {
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	migration, _ := migrate.NewWithDatabaseInstance(
		"file://../../../migration/mysql",
		"mysql",
		driver,
	)

	err := migration.Up()
	if err == nil {
		return nil
	}

	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

type InsertDummyFileParam struct {
	UniqueId  string
	Name      string
	Path      string
	Mimetype  string
	Extension string
	Size      int64
	CreatedAt int64
	UpdatedAt int64
}

func InsertDummyFile(db *sql.DB, p InsertDummyFileParam) error {
	query := "INSERT INTO file (id, name, path, mimetype, extension, size, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%d', '%d')"
	query = fmt.Sprintf(query, p.UniqueId, p.Name, p.Path, p.Mimetype, p.Extension, p.Size, p.CreatedAt, p.UpdatedAt)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
