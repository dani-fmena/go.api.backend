package models

import "time"

// go-pg naming convention (https://pg.uptrace.dev/models/)

// Book is the database table for holding the books
type Book struct {
	Id        uint		`example:"24"`
	Name      string    `pg:",unique" example:"The Book of Eli"`
	Items     uint      `pg:"default:0" example:"46"`
	CreatedAt time.Time `pg:"default:now()" example:"2021-03-12T02:11:03.292442-05:00"`
	UpdatedAt time.Time	`example:"0001-01-01T00:00:00Z"`
}

// TIP An model / entity can be an object with methods.
// The entities / models could be used by many different applications in the enterprise.
// This hst to encapsulate Enterprise wide business rules. Eg. Entity field transformation
