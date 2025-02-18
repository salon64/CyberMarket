package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	// "log"
	"net/http"
)

type Item struct {
	ItemID   int
	TypeID   int
	ItemName string
	// these are string pointers since ItemDescription and ImgURL can be null,
	// when scan is done null is converted to nil pointers
	ItemDescription *string
	ImgURL          *string
}

type SimpleItem struct {
	
}

// this functions writes out a json of all the items belonging to the users which is given in the url
// the data returned is described by the item struct
func createItem(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var newItem Item

	if r.Body == nil {
		log.Print("body was nil")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newItem)

	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "Error parsing json: %s", err.Error())
		log.Printf("Error parsing json: %s", err.Error())
		return
	}


	// write error and exit if json fails
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}
func listUserItems(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	// row, err := db.Query("SELECT ItemID FROM Inventory WHERE UserID = ? ODER BY ItemID", r.PathValue("id"))
	row, err := db.Query(`
		SELECT Inventory.ItemID, Inventory.TypeID, ItemTypes.ItemName, ItemTypes.ItemDescription, ItemTypes.ImgURL
		FROM Inventory
		INNER JOIN ItemTypes on Inventory.TypeID = ItemTypes.TypeID
		where Inventory.UserID = ?;`,
		r.PathValue("id"))

	if isErrLog(w, err) {
		return
	}
	// close the read data when the function
	defer row.Close()

	// tempora storage of items before sent to client
	var items []Item

	// for each row in our result
	for row.Next() {
		var item Item

		// read the columns from the row
		err := row.Scan(&item.ItemID, &item.TypeID, &item.ItemName, &item.ItemDescription, &item.ImgURL)

		// write error and exit if scan fails
		if err != nil {
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(*w, err.Error())
			return
		}

		items = append(items, item)
	}

	// to json
	json, err := json.MarshalIndent(items, "", "    ")

	// if conversion to json failed
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}
