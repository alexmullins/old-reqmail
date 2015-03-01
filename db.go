package main

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	createSettingsTableSchema = `CREATE TABLE IF NOT EXISTS settings (
			key string PRIMARY KEY,
			value string
		);`
	createEmailsTableSchema = `CREATE TABLE IF NOT EXISTS emails (
			buyer_name string PRIMARY KEY, 
			buyer_email string,
			active integer
		);`
	createReqTableSchema = `CREATE TABLE IF NOT EXISTS %s (
			req_no string,
			line_no integer,
			release_no integer,
			buyer_name string,
			PRIMARY KEY (req_no, line_no, release_no, buyer_name)
		);`

	subscribeBuyerStatement = `INSERT INTO emails (buyer_name, buyer_email, active) 
			VALUES (?, ?, ?);`

	unsubscribeBuyerStatement = `UPDATE emails
		SET active = 0
		WHERE buyer_name = ?`

	getAllBuyersStatement = `SELECT * FROM emails;`
)

// DB in charge of keeping track
// of A and B tables, email distribution list,
// and settings (current table). Also can get a collection of
// new requisitions.
type AppDB struct {
	*sqlx.DB
}

// Create the apps db with the correct structure.
func NewAppDB(name string) (*AppDB, error) {
	db, err := sqlx.Open("sqlite3", "app.db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Check if the db has the correct structure, if so return the db.
	// Must have the following tables:
	// 	1. settings
	//	2. emails
	//	3. a
	//	4. b
	_, err = db.Exec(createSettingsTableSchema)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createEmailsTableSchema)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(fmt.Sprintf(createReqTableSchema, "a"))
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(fmt.Sprintf(createReqTableSchema, "b"))
	if err != nil {
		return nil, err
	}
	return &AppDB{db}, nil
}

// Subscribe a buyer to receive updates
func (a *AppDB) SubscribeBuyer(b Buyer) (sql.Result, error) {
	r, err := a.Exec(subscribeBuyerStatement, b.Name, b.Email, true)
	return r, err
}

// Unsubscribe buyer from receiving updates
func (a *AppDB) UnsubscribeBuyer(b Buyer) (sql.Result, error) {
	r, err := a.Exec(unsubscribeBuyerStatement, b.Name)
	return r, err
}

// Get a list of Buyers from database
func (a *AppDB) GetBuyers() ([]Buyer, error) {
	buyers := make([]Buyer, 0)

	rows, err := a.Query(getAllBuyersStatement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		var email string
		var active bool

		rows.Scan(&name, &email, &active)
		b := Buyer{Name: name, Email: email, Active: active}
		buyers = append(buyers, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return buyers, nil
}

// Update the "new" requisition table
func (a *AppDB) UpdateReqTable(reqs []ReqLine) error {

}

// Get a list of updates (name -> [#1 req, #2 req, #3 req])
func (a *AppDB) GetUpdates() ([]ReqLine, error) {

}

// Simple Buyer
type Buyer struct {
	Name   string
	Email  string
	Active bool
}
