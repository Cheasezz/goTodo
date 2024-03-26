package app

import (
	"errors"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/source/github"
	"github.com/sirupsen/logrus"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func dbMigrate() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		logrus.Fatal("migrate: environment variable not declared: PG_URL")
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://schema", databaseURL)
		if err == nil {
			break
		}

		logrus.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		logrus.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logrus.Print("Migrate: no change")
		return
	}

	logrus.Printf("Migrate: up success")
}
