package main

import (
	"database/sql"
	. "http-proxy/internal"
	"log"
)

func main() {
	repeater := InitRepeater()
	repeater.DB, _ = sql.Open("sqlite3", "./data/data.db")
	defer repeater.DB.Close()

	log.Fatal(repeater.Server.ListenAndServe())
}
