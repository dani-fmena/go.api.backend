package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rubenv/sql-migrate"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"go.api.backend/schema/models"
	"go.api.backend/service/utils"
)

func Bootstrap(c *utils.SvcConfig) *pg.DB {

	// Creating the database connection and getting the context
	pgdb := pg.Connect(&pg.Options{Addr: c.Addr, User: c.User, Password: c.Pass, Database: c.Database})
	ctx := context.Background()

	// Check
	if err := pgdb.Ping(ctx); err != nil { panic(err) }

	return pgdb
}

// CreateSchema creates database schema for User and Story models.
//
// - db [*pg.DB] ~ Postgres database instance
//
// - testing [bool] ~ Tells if wee need to run the method just for checking purpose (everything OK with the table creation flow)
func CreateSchema(db *pg.DB, testing bool) {
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
}

// MkMigrations run the last migrations in a set.
func MkMigrations(c *utils.SvcConfig) {

	// For migrations we are using https://github.com/rubenv/sql-migrate. There is alternatives like
	// https://github.com/go-pg/migrations or https://github.com/golang-migrate/. ❗ An important thing to note is
	// that we are not use migration for the initial tables / schemas creation. For that matter we have go-pg
	// CreateSchema method. We can can create an endpoint http handler to create migration Up, Down
	// migrations methods (see https://github.com/rubenv/sql-migrate#usage)

	// Making db connection. ❗ Notice that we use database/sql because migration packages use it.
	// So we can't use go-pg connection instance for talk with the database
	pgCnxInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Pass, c.Database)

	pgdb, e := sql.Open("postgres", pgCnxInfo)
	if e != nil { panic(e) }

	// Making sure db disconnection just before method return / exit
	defer pgdb.Close()

	// Check
	if err := pgdb.Ping(); err != nil { panic(err) }

	// So far so good, setup migration source
	migrations := &migrate.FileMigrationSource{
		Dir: c.MigrationDir,
	}

	// Run the migrations
	n, err := migrate.Exec(pgdb, "postgres", migrations, migrate.Up)
	if err != nil { panic(err) }

	fmt.Printf("Applied %d migrations!\n", n)

	// Note that n can be greater than 0 even if there is an error: any migration that succeeded will remain
	// applied even if a later one fails.
}