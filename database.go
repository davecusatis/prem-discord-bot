package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// User represents a premium user
type User struct {
	email      string
	discordTag string
	startDate  int64
	endDate    int64
}

// Database is the wrapper around postgres connection
type Database struct {
	db *sql.DB
}

func newDatabase(connectionStr string) (*Database, error) {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Error creating database connection")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}
	return &Database{
		db: db,
	}, nil
}

func (db *Database) addUser(user *User) error {
	var email string
	err := db.db.QueryRow(fmt.Sprintf(`INSERT INTO users(Email, Discord, StartDate, EndDate ) VALUES('%s', '%s', '%d', '%d') RETURNING user_id`,
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
