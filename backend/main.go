package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

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

	// function that returns an array containing PubUser struct {id,name} in json format
	http.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		// allow CORS
		enableCors(&w)

		// execute sql query to get username id pairs
		row, err := db.Query("select userid, username from Users")

		// write error if the query returned error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		// close the row connection when function exits
		defer row.Close()

		// where to store result
		// note this was chosen instead of printing each row after read
		// this allows for retuning error if any row parsing fails.
		// the con of this that the result is buffered, which leads to an memory overhead
		var pubUsers []PubUser

		// prepare for next read
		for row.Next() {
			var user PubUser
			// read data into user struct
			err := row.Scan(&user.Id, &user.Name)

			// write error and exit if scan fails
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err.Error())
				return
			}
			// push user to the array
			pubUsers = append(pubUsers, user)
		}

		// convert to json
		// using MarshalIndent to make result pretty for debugging
		json, err := json.MarshalIndent(pubUsers, "", "    ")

		// write error and exit if json fails
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		// send json
		fmt.Fprint(w, string(json))
	})

	http.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// parse the form and return if err
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Malformed request %s", err.Error())
			return
		}

		// read pswd and name if either is empty return error
		pswd := r.Form.Get("pswd")
		username := r.Form.Get("name")
		if pswd == "" || username == "" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Malformed request")
			return
		}

		// add the user to Users
		res, err := db.Exec("insert into Users(Username,`Password(Hash)`,Wallet, role) values (?, ?, 0, 0);", username, pswd)
		// if error write error and exit
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
			return
		}
		// get the id if the row
		id, err := res.LastInsertId()
		// if error write error and exit
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
			return
		}

		// return the id
		fmt.Fprintf(w, "%d", id)
	})

	http.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
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
