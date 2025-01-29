package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AddUser struct {
	Name string `json:"name,omitempty"`
	Pswd string `json:"pswd,omitempty"`
}

type PubUser struct {
	Id   int
	Name string
}

func listAllUsers(w *http.ResponseWriter, _ *http.Request, db *sql.DB) {
	// execute sql query to get username id pairs
	row, err := db.Query("select userid, username from Users")

	// write error if the query returned error
	if isErrLog(w, err) {
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
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(*w, err.Error())
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
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}

func addUser(w *http.ResponseWriter, r *http.Request, db *sql.DB) {

	var useradd AddUser

	if r.Body == nil {
		log.Print("body was nil")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&useradd)

	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "Error parsing json: %s", err.Error())
		log.Printf("Error parsing json: %s", err.Error())
		return
	}

	// test if values aer empty strings
	if useradd.Name == "" || useradd.Pswd == "" {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, "Either Name or Pswd is an empty string")
		log.Print("Either Name or Pswd is an empty string")
		return
	}

	// add the user to Users
	res, err := db.Exec("insert into Users(Username,`Password(Hash)`,Wallet, role) values (?, ?, 0, 0);", useradd.Name, useradd.Pswd)
	// if error write error and exit
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("add user error: %s", err)
		fmt.Fprintln(*w, err.Error())
		return
	}
	// get the id if the row
	id, err := res.LastInsertId()
	// if error write error and exit
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(*w, err.Error())
		return
	}

	// return the id
	fmt.Fprintf(*w, "%d", id)
}
