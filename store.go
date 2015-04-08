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

// Interface to define storage for requisitions
type ReqRepository interface {
	SubscribeBuyer(*Buyer) error
	UnsubscribeBuyer(*Buyer) error
	UpdateReport(*ReqReport) error
	GetNewReqs() (*NewReqReport, error)
}

// DB in charge of keeping track
// of A and B tables, email distribution list,
// and settings (current table). Also can get a collection of
// new requisitions. Concrete implementation of ReqRepository interface
type SqliteRepo struct {
	*sqlx.DB
}

// Create the apps db with the correct structure.
func NewSqliteRepo(name string) (*SqliteRepo, error) {
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
	if err = createDBStructure(db); err != nil {
		return nil, err
	}

	return &SqliteRepo{db}, nil
}

func createDBStructure(db *sqlx.DB) error {
	_, err := db.Exec(createSettingsTableSchema)
	if err != nil {
		return err
	}
	_, err = db.Exec(createEmailsTableSchema)
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf(createReqTableSchema, "a"))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf(createReqTableSchema, "b"))
	if err != nil {
		return err
	}
	return nil
}

// Subscribe a buyer to receive updates
func (a *SqliteRepo) SubscribeBuyer(b *Buyer) (sql.Result, error) {
	r, err := a.Exec(subscribeBuyerStatement, b.Name, b.Email, true)
	return r, err
}

// Unsubscribe buyer from receiving updates
func (a *SqliteRepo) UnsubscribeBuyer(b *Buyer) (sql.Result, error) {
	r, err := a.Exec(unsubscribeBuyerStatement, b.Name)
	return r, err
}

func (s *SqliteRepo) UpdateReport(r *ReqReport) error {
	// TODO
	return nil
}

func (s *SqliteRepo) GetNewReqs() *NewReqReport {
	// TODO
	return nil
}

// Struct to hold the data of "New Reqs"
type NewReq struct {
	ReqNo     string
	LineNo    string
	ReleaseNo string
	BuyerName string
	Email     string
}

type NewReqReport []NewReq

// Simple Buyer
type Buyer struct {
	Name   string
	Email  string
	Active bool
}
