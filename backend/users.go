package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// struct used for adding a user or updating their info
type AddUser struct {
	Name string `json:"name,omitempty"`
	Pswd string `json:"pswd,omitempty"`
}

// a struct containing the fields that are public for a user
type PubUser struct {
	Id   int
	Name string
}

type UserMoney struct {
	ID int
	Money int
}

// list all the users
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

func updateUserInfo(w *http.ResponseWriter, r *http.Request, db *sql.DB) {

	var data map[string]string

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)

	// TODO do token check
	token := data["token"]
	if token == "" {
		log.Print("failed to authenticate token")
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, "failed to authenticate token")
		return
	}

	new_name, name_ok := data["new_name"]
	new_pswd, pswd_ok := data["new_pswd"]

	if name_ok {
		db.Exec("UPDATE Users SET Username=? WHERE UserID = ?", new_name, r.PathValue("id"))
	}
	if pswd_ok {
		db.Exec("UPDATE Users SET `Password(Hash)`=? WHERE UserID = ?", new_pswd, r.PathValue("id"))
	}

	// TODO return old name and pswd

	(*w).WriteHeader(http.StatusOK)
}

func addMoneyToUser(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data UserMoney
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)

	// TODO token stuff
	db.Exec("UPDATE Users SET Wallet = Wallet + ? WHERE UserID = ?", data.Money, data.ID)

	(*w).WriteHeader(http.StatusOK)
}
