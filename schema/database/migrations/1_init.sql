-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
-- `create table sample` -> CREATE TABLE IF NOT EXISTS table_name ();
CREATE UNIQUE INDEX book_name_idx ON books (lower (name));          -- Its also good practice create an index for relationships


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE people;
