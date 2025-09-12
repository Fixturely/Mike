package utils

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"mike/config"
	"runtime"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	instance *bun.DB
	dbError  error
)

// GetDatabase returns a singleton database connection
func GetDatabase() *bun.DB {
	return instance
}

func init() {
	log.Println("initializing database")
	instance, dbError = initDatabase()
	if dbError != nil {
		log.Fatalf("error initializing database: %v", dbError)
	}
}

func initDatabase() (*bun.DB, error) {
	cfg := config.GetConfig()
	databaseCfg := cfg.Database
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	opts := []pgdriver.Option{
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", databaseCfg.Host, databaseCfg.Port)),
		pgdriver.WithDatabase(databaseCfg.Name),
		pgdriver.WithUser(databaseCfg.User),
		pgdriver.WithPassword(databaseCfg.Password),
		pgdriver.WithTimeout(15 * time.Second),
		pgdriver.WithDialTimeout(15 * time.Second),
	}
	if databaseCfg.SSLMode {
		// TODO: Revisit the guardrails ignore.
		opts = append(opts, pgdriver.WithTLSConfig(&tls.Config{ServerName: databaseCfg.Host})) // guardrails-disable-line
	} else {
		opts = append(opts, pgdriver.WithInsecure(true))
	}
	pgconn := pgdriver.NewConnector(opts...)

	sqldb := sql.OpenDB(pgconn)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	db := bun.NewDB(sqldb, pgdialect.New())

	// Show queries in logs for development
	if cfg.Environment == "development" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// Ensure the database can connect.
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatalf("error initializing database, unable to SELECT 1: %v", err)
	}

	return db, nil
}
