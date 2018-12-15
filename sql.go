// Package sql provides a GoHub database implementation backed by a SQL db.
package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/mughub/mughub/db"
	"github.com/spf13/viper"
	"io"
)

// DB represents a general SQL database and is built on top of the
// database/sql and database/sql/driver packages from Go's standard library.
//
type DB struct {
	drvr   *sql.DB
	schema graphql.Schema
}

// Init checks basic config values and then calls sql.Open() to
// establish a connection to your SQL database.
//
func (d *DB) Init(schema io.Reader, cfg *viper.Viper) (err error) {
	driver := cfg.GetString("name")
	if driver == "" {
		return db.ErrUnknownDB{Name: "empty database name provided"}
	}

	err = db.ErrUnknownDB{Name: driver}
	for _, drvr := range sql.Drivers() {
		if drvr == driver {
			err = nil
		}
	}
	if err != nil {
		return
	}

	dsn := cfg.GetString("dsn")
	if dsn == "" {
		return errors.New("db: dsn must be non-empty")
	}

	d.drvr, err = sql.Open(driver, dsn)
	if err != nil {
		return
	}

	// TODO: Add SQL Resolvers to GraphQL API Schema
	return
}

// Do executes a GraphQL request as a SQL request.
func (d *DB) Do(ctx context.Context, req string, vars map[string]interface{}) *db.Result {
	res := graphql.Do(graphql.Params{
		Context:        ctx,
		RequestString:  req,
		VariableValues: vars,
	})
	return &db.Result{
		Data: res.Data,
		// TODO: Transfer errors
	}
}
