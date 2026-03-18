package database

import (
	"fmt"

	"github.com/forbole/juno/v5/database"
	"github.com/forbole/juno/v5/database/postgresql"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ database.Database = &DB{}

// DB represents a PostgreSQL database with expanded features.
// so that it can properly store custom BigDipper-related data.
type DB struct {
	G *gorm.DB
	*postgresql.Database
	Sqlx *sqlx.DB
}

// Builder allows to create a new Db instance implementing the db.Builder type
func Builder(ctx *database.Context) (database.Database, error) {
	database, err := postgresql.Builder(ctx)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid configuration database, must be PostgreSQL")
	}
	gdb, err := gorm.Open(postgres.Open(ctx.Cfg.URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error while opening the connection to the database: %s", err)
	}
	return &DB{
		G:        gdb,
		Database: psqlDb,
		Sqlx:     sqlx.NewDb(psqlDb.SQL.DB, "postgresql"),
	}, nil
}

// Cast allows to cast the given db to a Db instance
func Cast(db database.Database) *DB {
	bdDatabase, ok := db.(*DB)
	if !ok {
		panic(fmt.Errorf("given database instance is not a Db"))
	}
	return bdDatabase
}

// SQLStatement holds a single SQL statement and its bind vars for ordered execution.
type SQLStatement struct {
	SQL  string
	Vars []interface{}
}

func (db *DB) ExecuteStatements(statements []SQLStatement) error {
	if len(statements) == 0 {
		return nil
	}
	tx := db.G.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, stmt := range statements {
		if err := tx.Exec(stmt.SQL, stmt.Vars...).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

