package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Interface defining a source of open requisition report
type ReqSource interface {
	GetReqReport() (*ReqReport, error)
}

// Concrete implementation of ReqSource interface
type IFS struct {
}

func (i *IFS) GetReqReport() (*ReqReport, error) {
	return nil, nil
}

type ReqLine struct {
	ReqNo     string
	LineNo    int
	ReleaseNo int
	BuyerName string
}

type ReqReport []ReqLine

// Req Source from .csv file
// Will watch a directory for a new .csv file
// with a filename incremented by +1 for new "update"
type CSVSoure struct {
	dir    string     // Dir to watch
	report *ReqReport // lazy load the req report from the current file when calling GetReqReport()

	mu      sync.Mutex // protects everything beneath
	current int        // First file must start at "1.csv"
	updated bool       // Do we have an updated file since the last GetReqReport() call
}

// Watch a dir for a new file
// Polling ever *duration* seconds
// Should be run in its own goroutine
func (c *CSVSoure) Watch(dir string) {
	// Set current to the "highest" .csv file in the dir
	// Get a list of filename in the current directory minus the extension
	// Convert list to integer and sort highest to lowest getting element at index 0
	// Set current string to highest filename
	c.dir = dir

	if c.current == 0 {
		id := c.getHighestFilename()
		c.updateCurrent(id)
	}

	// Start a for loop waiting every duration seconds polling for a higher id file
	var countdown = time.NewTicker(15 * time.Second)

	for {
		select {
		case <-countdown.C:
			id := c.getHighestFilename()
			if id > c.current {
				c.updateCurrent(id)
			}
		}
	}
}

func (c *CSVSoure) updateCurrent(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = id
	c.updated = true
	log.Printf("Set c.current to: %d", c.current)
}

func (c *CSVSoure) getHighestFilename() int {
	files, err := ioutil.ReadDir(c.dir)
	if err != nil {
		log.Fatalf("Failed to read dir: %s", err)
	}
	// Since we know this is first pull we will get the last index
	last := files[len(files)-1].Name()
	parts := strings.Split(last, ".")
	name := parts[0]
	ext := parts[1]

	// Error out if we aren't .csv files
	if ext != "csv" {
		log.Fatalf("CSVSource can only accept .csv files")
	}

	// Error out if we can't convert name to integer
	id, err := strconv.Atoi(name)
	if err != nil {
		log.Fatal("Must only have numerically incremented .csv files")
	}
	return id
}

func (c *CSVSoure) GetReqReport() (*ReqReport, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// We've updated since we've called GetReqReport
	if c.updated {
		report, err := c.parseReqReport()
		if err != nil {
			return nil, err
		}
		c.updated = false
		c.report = report
		return c.report, nil
	}
	return c.report, nil
}

func (c *CSVSoure) parseReqReport() (*ReqReport, error) {
	return nil, errors.New("Dumb Error")
}
