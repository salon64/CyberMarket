package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	// "log"
	"net/http"
)

type ItemTypeInformation struct {
	ItemName string
	ItemDescription string 
	ImgURL string 
	ShortDescription string
}

// adds a new itemtype to the ItemType table, 
// TODO implement description and img handling 
func createNewItemType(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data ItemTypeInformation

	if r.Body == nil {
		log.Print("body was nil")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)

	// TODO ERROR handling
	t, _ := db.Begin()
	res, err := t.Exec("insert into ItemTypes (ItemName, ItemDescription, ImgURL, ShortDescription) values (?, null, null, ?);", data.ItemName, data.ShortDescription)

	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("add listing error: %s", err)
		fmt.Fprintf(*w, "add listing error: %s", err.Error())
		return
	}
	// TODO: WE NEED A TRANSACTION HERE
	TypeID, err := res.LastInsertId()
	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(*w, err.Error())
		return
	}

	// return the TypeId
	t.Commit()
	fmt.Fprintf(*w, "%d", TypeID)
}
