package database

import (
	"database/sql"
	"fmt"

	"github.com/apm-dev/oha/src/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func NewConnection(config *config.Config) (*sql.DB, error) {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DB,
	)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ApplyDatabaseMigrations(config *config.Config) {
	log.Infof("going to apply database migrations")
	databaseURL := fmt.Sprintf(
		"postgresql://%s:%d/%s?user=%s&password=%s&sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.DB, config.Database.User, config.Database.Password,
	)
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("can't init migration '%s'", err)
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatalf("unable to apply migrations '%s'", err)
	}
	_, _ = m.Close()
	log.Infof("all database migrations has been successfully applied")
}
