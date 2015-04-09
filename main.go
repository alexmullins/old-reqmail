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

	// Create datasource
	csv := &CSVSoure{}
	go csv.Watch("./source/")

	// // Create datastore
	// store := &MemoryStore{}

	// // Create email updater
	// // TODO: pass in username, password, and email server address
	// emailer := &Emailer{}

	log.Println("Created timer: ", *pollInterval)

	app := &App{
		Poll:   pollInterval,
		Source: csv,
		// Store:   store,
		// Updater: emailer,
	}

	app.Run()

}

type App struct {
	Poll   *time.Duration // poll period to data source
	Source ReqSource      // source to pull ReqReport
}

func (a *App) Run() {
	countdown := time.NewTicker(*a.Poll)

	for {
		select {
		case <-countdown.C:
			// Signaled to poll source for req data
			log.Println("Got signal to poll IFS.")
			report := a.pullReqReport()
			log.Printf("report: %v", report)

			// // Pass the new req report to datastore to find updates
			// var store ReqRepository
			// store.UpdateReport(report)

			// // Get updates
			// updates, err := store.GetNewReqs()
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// // Pass updates to updater
			// var updater Updater
			// updater.SendUpdates(updates)
		}
	}
}

func (a *App) pullReqReport() *ReqReport {
	report, err := a.Source.GetReqReport()
	if err != nil {
		log.Fatal(err)
	}
	return report
}
