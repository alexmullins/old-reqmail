package main

import "testing"

func TestCSVSourceWatchReadsHighestIdFilename(t *testing.T) {
	csv := CSVSoure{}
	csv.watch("./source/")

}
