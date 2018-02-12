package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

func (db *Database) userExists(email string) (bool, error) {
	var exists bool
	existsQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM users WHERE email = '%s')", email)
	err := db.db.QueryRow(existsQuery).Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, nil
}

func (db *Database) addOrUpdateUser(user *User) error {
	exists, err := db.userExists(user.email)
	if err != nil {
		return err
	}

	var query string
	if exists {
		query = fmt.Sprintf(`UPDATE users SET product = '%s', discord = '%s', start_date = '%d', end_date = '%d' WHERE email = '%s' RETURNING id`,
			user.product,
			user.discordTag,
			user.startDate,
			user.endDate,
			user.email)
	} else {
		query = fmt.Sprintf(`INSERT INTO users(email, product, discord, start_date, end_date ) VALUES('%s', '%s', '%s', '%d', '%d') RETURNING id`,
			user.email,
			user.product,
			user.discordTag,
			user.startDate,
			user.endDate)
	}

	var ID string
	err = db.db.QueryRow(query).Scan(&ID)
	if err != sql.ErrNoRows && err != nil {
		return err
	}
	return nil
}

func (db *Database) getDurationByProductID(productID string) (*time.Duration, error) {
	var duration string
	err := db.db.QueryRow(fmt.Sprintf("SELECT duration FROM products WHERE product_id = '%s'", productID)).Scan(&duration)

	if err != sql.ErrNoRows && err != nil {
		return nil, err
	}
	ret, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
