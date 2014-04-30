package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// TODO(benbjohnson): Read STDIN into the database.

	// Serve root handler.
	fmt.Printf("Listening on http://localhost%s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, sieve.NewHandler(db)))
}
