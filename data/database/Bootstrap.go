package database

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go.api.backend/data/models"
	"go.api.backend/services"
)

// Bootstrap create the database engine object
//
// - c [*pg.DB] ~ Pointer to the database connection
func Bootstrap(c *services.SvcConfig) *pg.DB {

	// Creating the database connection and getting the context
	pgdb := pg.Connect(&pg.Options{Addr: c.Host, User: c.User, Password: c.Pass, Database: c.Database})
	ctx := context.Background()

	// Checking
	if err := pgdb.Ping(ctx); err != nil {
		panic(err)
	}

	return pgdb
}

// CreateSchema creates database schema for User and Story models.
//
// - db [*pg.DB] ~ Postgres database instance
//
// - testing [bool] ~ Tells if wee need to run the method just for checking purpose (everything OK with the table creation flow)
func CreateSchema(db *pg.DB, testing bool) error {
	schemas := []interface{} {
		(*models.Book)(nil),
	}

	for _, model := range schemas {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions {
			IfNotExists: true,							// If the table exist don't create it
			Temp:        testing,						// Don commit the change to the database
		})

		if err != nil && !testing {
			println(err.Error())
		}
	}
	return nil
}
