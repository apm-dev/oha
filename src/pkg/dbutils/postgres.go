package dbutils

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type PostgresConfig struct {
	Host   string
	Port   uint
	User   string
	Passwd string
	DB     string
	// MigrationPathURL in form of file://migrations
	MigrationPathURL string
}

func NewPgsqlConnection(config *PostgresConfig) (*sql.DB, error) {
	pgsqlConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Passwd, config.DB,
	)
	db, err := sql.Open("postgres", pgsqlConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ApplyPgsqlDatabaseMigrations(config *PostgresConfig) {
	log.Infof("going to apply database migrations")
	databaseURL := fmt.Sprintf(
		"postgresql://%s:%d/%s?user=%s&password=%s&sslmode=disable",
		config.Host, config.Port, config.DB, config.User, config.Passwd,
	)
	m, err := migrate.New(config.MigrationPathURL, databaseURL)
	if err != nil {
		log.Fatalf("can't init migration '%v'", err)
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatalf("unable to apply migrations '%v'", err)
	}
	_, _ = m.Close()
	log.Infof("all database migrations has been successfully applied")
}
