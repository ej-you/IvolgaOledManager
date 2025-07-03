// The i2c_scanner binary scans all I2C devices in the /dev directory and
// prints out busy I2C bus addresses.
// To check that I2C address is busy scanner sends zero byte to it and
// waits for response without errors.
package main

import (
	"flag"
	"io"
	"log"
	"os"

	"IvolgaOledManager/internal/pkg/i2c"
)

func main() {
	verboseFlag := flag.Bool("v", false, "verbose output")
	flag.Parse()

	log.Println("Start scanning i2c...")
	if err := runI2CScanner(*verboseFlag); err != nil {
		log.Fatal(err)
	}
}

func runI2CScanner(verboseFlag bool) error {
	logWriter := io.Discard
	if verboseFlag {
		logWriter = os.Stdout
	}
	// init scanner
	scanner := i2c.NewScanner(logWriter)
	// scan all i2c buses ans addresses
	scanResult, err := scanner.ScanAll()
	if err != nil {
		return err
	}
	// print out found results
	log.Println(scanResult)
	return nil
}
