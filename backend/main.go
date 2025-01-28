package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {

	cfg := mysql.Config{
		User:      os.Getenv("DBUSER"),
		Passwd:    os.Getenv("DBPASS"),
		Net:       "tcp",
		ParseTime: true,
		Addr:      os.Getenv("DBHOST"),
		DBName:    "main_db", // this is the name that is defined
	}
	var err error // used here so that it makes db global

	// try to connect 30 times each waiting 2 seconds
	log.Printf("Connecting to database at host %s, with DSN %s", os.Getenv("DBHOST"), cfg.FormatDSN())
	for i := 1; i < 31; i++ {

		// try to connect with cfg as connection options
		db, err = sql.Open("mysql", cfg.FormatDSN())

		// if connection attempt returned error log
		if err != nil {
			log.Printf("Error connecting, attempt %d/30. %s. ", i, err.Error())
			time.Sleep(time.Second * 2)
		} else {
			log.Print("Connection established")
			break
		}
	}

	// if error still exist after 30 connection tries, abort
	if err != nil {
		log.Fatal("aborting failed to establish connection")
	}
	defer db.Close() // close db when main exits

	// try to ping 30 times each waiting 2 seconds

	var pingErr error // test
	for i := 1; i < 31; i++ {
		// pinging db
		pingErr = db.Ping()

		// if ping
		if pingErr != nil {
			log.Printf("Error pinging, attempt %d/30. %s. ", i, pingErr.Error())
			time.Sleep(time.Second * 2)
		} else {
			log.Print("Connection verified")
			break
		}
	}
	// if error still exist after 30 ping attempts, abort
	if pingErr != nil {
		log.Fatal("aborting failed to ping")
	}

	// TODO set max life time and other
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s ", os.Getenv("DBUSER"))
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		row, err := db.Query("SHOW TABLES;")

		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			for row.Next() {
				var table string
				err := row.Scan(&table)
				if err != nil {
					fmt.Fprintln(w, err.Error())
				} else {
					fmt.Fprintln(w, table)
				}
			}
		}
	})

	http.ListenAndServe(":80", nil)
}
