package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	emailAddr     = flag.String("email", "user@email.com", "Sender's email account.")
	emailPassword = flag.String("password", "", "Sender's password.")

	createSettingsTableSchema = `CREATE TABLE IF NOT EXISTS settings (
			key string,
			value string
		);`
	createEmailsTableSchema = `CREATE TABLE IF NOT EXISTS emails (
			buyer_name string,
			email string
		);`
	createATableSchema = `CREATE TABLE IF NOT EXISTS a (
			req_no string,
			line_no integer,
			release_no integer,
			buyer_name string
		);`
	createBTableSchema = `CREATE TABLE IF NOT EXISTS b (
			req_no string,
			line_no integer,
			release_no integer,
			buyer_name string
		);`
)

// Create the apps db with the correct structure.
func createAppDB(name string) (*sqlx.DB, error) {
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
	_, err = db.Exec(createATableSchema)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createBTableSchema)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	flag.Parse()

	db, err := createAppDB("app.db")
	if err != nil {
		log.Fatalln("couldn't create app db: %v", err)
	}

	result := db.MustExec(`INSERT INTO emails (buyer_name, email) VALUES ("JOSEPHMULLINS", "jamminalex@gmail.com");`)
	id, _ := result.LastInsertId()
	fmt.Printf("last inserted id: %v\n", id)
	rows, _ := db.Query("SELECT * FROM emails;")
	for rows.Next() {
		var name string
		var email string
		rows.Scan(&name, &email)
		fmt.Printf("name: %v\n", name)
		fmt.Printf("email: %v\n", email)
	}
}
