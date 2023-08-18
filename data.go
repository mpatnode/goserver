package main

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Addr is the User address structure
type Addr struct {
	Line1       string `json:"line_1"`
	Line2       string `json:"line_2,omitempty"`
	City        string `json:"city"`
	Subdivision string `json:"subdivision"`
	PostalCode  string `json:"postal_code"`
}

// User represents a user table row
type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name,omitempty"`
	LastName    string `json:"last_name"`
	Address     Addr   `json:"address"`
	StripeID    string `json:"stripe_id"`
	BraintreeID string `json:"bt_id"`
	CreatedAt   int64  `json:"timestamp"`
}

var _db *sql.DB = nil

func getDB() *sql.DB {
	if _db == nil {
		var err error
		_db, err = sql.Open("sqlite3", "database/fast.db")
		if err == nil {
			_, err = _db.Exec("CREATE TABLE IF NOT EXISTS users (id text UNIQUE, stripeID text, btID text, email text UNIQUE, firstName text, middleName text, lastName text, created integer)")
			_, err = _db.Exec("CREATE TABLE IF NOT EXISTS addresses (uid text, line1 text, line2 text, city text, subdivision text, postalCode text)")
		}

		if err != nil {
			log.Fatal(err)
		}
	}
	return _db
}

// DBAddUser to the database
func DBAddUser(u *User) error {
	var db = getDB()
	// Just being lazy here.  Would typically let the DB generate the id.  Looked like you wanted a uuid
	id := uuid.New().String()

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO users (id, stripeID, btID, email, firstName, middleName, lastName, created) values (?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	_, err := stmt.Exec(id, u.StripeID, u.BraintreeID, u.Email, u.FirstName, u.MiddleName, u.LastName, time.Now().Unix())
	if err == nil {
		stmt, _ = tx.Prepare("INSERT INTO addresses (uid, line1, line2, city, subdivision, postalCode) values (?,?,?,?,?,?)")
		defer stmt.Close()
		_, err = stmt.Exec(id, u.Address.Line1, u.Address.Line2, u.Address.City, u.Address.Subdivision, u.Address.PostalCode)
		tx.Commit()
		u.ID = id
	} else {
		tx.Rollback()
	}
	return err
}

// DBGetUser from the database using the Fast ID or an email
func DBGetUser(key string) (*User, error) {
	var db = getDB()
	var u User
	query := "SELECT id, stripeID, btID, email, firstName, middleName, lastName, line1, line2, postalCode, subdivision, created from users JOIN addresses ON id = uid "
	if strings.Index(key, "@") > 0 {
		query += "WHERE email = ?"
	} else {
		query += "WHERE id = ?"
	}

	row := db.QueryRow(query, key)
	// XXX: Seems like there is a nice opportunity here to create an ORM which uses the Go object inspection to do this.
	// Something like the json tags?
	err := row.Scan(&u.ID, &u.StripeID, &u.BraintreeID, &u.Email, &u.FirstName, &u.MiddleName, &u.LastName, &u.Address.Line1, &u.Address.Line2, &u.Address.PostalCode, &u.Address.Subdivision, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// DBDeleteUser remove the user from the DB
func DBDeleteUser(key string) error {
	var db = getDB()
	u, err := DBGetUser(key)
	if err == nil && u != nil {
		tx, _ := db.Begin()
		db.Exec("DELETE FROM users WHERE id = ?", u.ID)
		db.Exec("DELETE FROM addresses WHERE uid = ?", u.ID)
		err = tx.Commit()
	}
	return err
}
