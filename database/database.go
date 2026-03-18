package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/zakiverse/zakiverse-api/config"
)

func InitDB(conf config.ConfigConstant, credential config.ConfigCredential) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		conf.Database.Host,
		conf.Database.Port,
		credential.DatabaseUsername,
		credential.DatabasePassword,
		conf.Database.Name,
		conf.Application.Timezone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping PostgreSQL: %v", err)
	}

	db.SetMaxIdleConns(conf.Database.MaxIdleConnection)
	db.SetMaxOpenConns(conf.Database.MaxOpenConnection)
	db.SetConnMaxLifetime(time.Duration(conf.Database.MaxConnectionLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(conf.Database.MaxConnectionIdleTime) * time.Minute)

	return db
}
