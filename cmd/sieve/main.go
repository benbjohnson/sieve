package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/benbjohnson/sieve"
)

var (
	addr = flag.String("addr", ":6900", "HTTP address")
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: sieve [opts]")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	// Read configuration.
	flag.Usage = usage
	flag.Parse()

	// Setup the database.
	var db = sieve.NewDB()

	// Read STDIN into the database.
	go load(os.Stdin, db)

	// Serve root handler.
	fmt.Printf("Listening on http://localhost%s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, sieve.NewHandler(db)))
}

// Parses a reader and streams it into the database.
func load(r io.Reader, db *sieve.DB) {
	var decoder = json.NewDecoder(r)
	for {
		// Parse individual row from JSON.
		var row = &sieve.Row{}
		if err := decoder.Decode(&row.Data); err == io.EOF {
			break
		} else if err != nil {
			log.Println("err:", err)
			continue
		}

		// Add it to the database.
		db.Append(row)
	}
}
