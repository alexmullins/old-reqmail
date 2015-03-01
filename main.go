package main

import (
	"flag"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	emailAddr     = flag.String("email", "user@email.com", "Sender's email account.")
	emailPassword = flag.String("password", "", "Sender's password.")
)

func main() {
	flag.Parse()

	db, err := NewAppDB("app.db")
	if err != nil {
		log.Fatalln("couldn't create app db: %v", err)
	}

	b := Buyer{Name: "ALEXMULLINS", Email: "jamminalex@gmail.com"}

	_, err = db.UnsubscribeBuyer(b)
	if err != nil {
		log.Printf("Couldn't unsubscribe buyer: %v", err)
	}

	buyers, err := db.GetBuyers()
	if err != nil {
		log.Printf("Couldn't get buyers, %v\n", err)
	}
	log.Println(buyers)
}
