package models

import "time"

// go-pg naming convention (https://pg.uptrace.dev/models/)

// Book is the database table for holding the books
type Book struct {
	Id        int64
	Name      string	`pg:",unique"`
	Items     uint64    `pg:"default:0"`
	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time
}

// TableName override the XORM naming convention for a table (https://gobook.io/read/gitea.com/xorm/manual-en-US/chapter-02/3.tags.html)
func (*Book) TableName() string {
	return "Book"
}