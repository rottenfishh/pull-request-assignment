package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
	DSN  string
}

func NewDb(ctx context.Context, dsn string) (*DB, error) {
	db := &DB{DSN: dsn}
	err := db.connectDB(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *DB) connectDB(ctx context.Context) error {

	var err error
	d.Pool, err = pgxpool.New(ctx, d.DSN)

	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := d.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database!")
	return nil
}

func (d *DB) RunMigrations() error {

	m, err := migrate.New("file://database/migrations", d.DSN)
	if err != nil {
		return fmt.Errorf("unable to run migrations: %v", err)
	}

	//err = m.Force(1)
	//if err != nil {
	//	return fmt.Errorf("failed to force migrations: %v", err)
	//}
	//if err != nil {
	//	return fmt.Errorf("failed to down migrations: %v", err)
	//}
	err = m.Down()
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
