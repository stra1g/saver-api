package database

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stra1g/saver-api/internal/infra/config"
)

var (
	db   *pgxpool.Pool
	once sync.Once
)

func NewPostgresDatabase(config *config.Config) *pgxpool.Pool {
	once.Do(func() {
		dbDSN := config.GetDatabaseDSN()

		poolConfig, err := pgxpool.ParseConfig(dbDSN)
		if err != nil {
			log.Fatalln("Unable to parse connection url:", err)
		}

		db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatalln("Unable to create connection pool:", err)
		}

		if err := db.Ping(context.Background()); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})

	return db
}
