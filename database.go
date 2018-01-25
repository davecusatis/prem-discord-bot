package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	email      string
	discordTag string
	startDate  string
	endDate    string
}

type Database struct {
	db *sql.DB
}

func newDatabase(connectionStr string) *sql.DB {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Error creating database connection")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}
	return db
}

func (db *Database) addUser(user *User) error {
	var email string
	err := db.db.QueryRow(
		fmt.Sprintf(`INSERT INTO users(Email, Discord, StartDate, EndDate ) VALUES('%s', '%s', '%s', '%s') RETURNING id`,
			user.email,
			user.discordTag,
			user.startDate,
			user.endDate,
		)).Scan(&email)

	if err != sql.ErrNoRows && err != nil {
		return err
	}
	log.Printf("Added user")
	return nil
}
