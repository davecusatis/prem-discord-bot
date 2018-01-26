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
	product    string
	discordTag string
	startDate  int64
	endDate    int64
}

// Database is the wrapper around postgres connection
type Database struct {
	db *sql.DB
}

func newDatabase() (*Database, error) {
	dbUser := mustGetConfigValue("DB_USER")
	dbPassword := mustGetConfigValue("DB_PASSWORD")
	dbPort := mustGetConfigValue("DB_PORT")
	dbName := mustGetConfigValue("DB_NAME")
	connectionStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s", dbUser, dbName, dbPassword, dbPort)
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
	err := db.db.QueryRow(fmt.Sprintf(`INSERT INTO users(email, product, discord, start_date, end_date ) VALUES('%s', '%s', '%s', '%d', '%d') RETURNING user_id`,
		user.email,
		user.product,
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
