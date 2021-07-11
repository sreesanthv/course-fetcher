package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const courseTable = `
CREATE TABLE courses (
id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
name text NOT NULL,
description text NOT NULL,
author text NOT NULL,
UNIQUE(name, author),
PRIMARY KEY (id)
)`

const courseIndex = "CREATE INDEX courses_idx ON courses (name, author, description)"

func init() {
	up := []string{
		courseTable,
		courseIndex,
	}

	down := []string{
		`DROP TABLE courses`,
		`DROP INDEX IF EXISTS courses_idx`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("up initial")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("down initial")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
