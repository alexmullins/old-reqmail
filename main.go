package main

import (
	"flag"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	emailAddr     = flag.String("email", "user@email.com", "Sender's email account.")
	emailPassword = flag.String("password", "", "Sender's password.")
	pollInterval  = flag.Duration("poll", 15*time.Second, "IFS polling interval.")
)

func main() {
	flag.Parse()

	// Create local sqlite database
	_, err := NewReqRepo("app.db")
	if err != nil {
		log.Fatalln("couldn't create app db: %v", err)
	}
	log.Println("Connected to database.")

	// // Connect to IFS via ODBC
	// ifs := NewIFSConn("connection string")
	// log.Println("Connected to IFS")

	// Create a timer that signals a
	// goroutine to start and pull an open req report
	// every x minutues from IFS.
	countdown := time.NewTicker(*pollInterval)
	log.Println("Created timer: ", *pollInterval)
	for {
		select {
		case <-countdown.C:
			log.Println("Got signal to poll IFS.")
			// go ifs.poll()
			return
		}
	}
}
