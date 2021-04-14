package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func Connect() (*Postgres, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", //os.Getenv("USERS_DB_HOST"),
		5432,
		"laurent",  //os.Getenv("USERS_DB_USERNAME"),
		"laurent",  //os.Getenv("USERS_DB_PASSWORD"),
		"macstats") //os.Getenv("USERS_DB_NAME"))
	dbpool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &Postgres{pool: dbpool}, nil
}

func (db *Postgres) InsertBattery(host string, stamp time.Time, bat BatteryInfo) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	prep := "INSERT INTO battery(host, stamp, metrics) VALUES ($1, $2, $3)"
	r, err := conn.Exec(context.Background(), prep,
		host,
		stamp,
		bat)
	if err != nil {
		return err
	}
	if r.RowsAffected() != 1 {
		return fmt.Errorf("insert into ssd failed")
	}
	return nil
}
func (db *Postgres) InsertSSD(host string, stamp time.Time, ssd SsdInfo) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	prep := "INSERT INTO ssd(host, stamp, metrics) VALUES ($1, $2, $3)"
	r, err := conn.Exec(context.Background(), prep,
		host,
		stamp,
		ssd)
	if err != nil {
		return err
	}
	if r.RowsAffected() != 1 {
		return fmt.Errorf("insert into ssd failed")
	}
	return nil
}
