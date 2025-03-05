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
	// these are pointers since ItemDescription, IsListed and ImgURL can be null,
	// when using scan null is converted to nil pointers
	IsListed        *int
	ItemDescription *string
	ImgURL          *string
}

type SimpleItem struct {
	UserID   int
	ItemType int
}

// this functions writes out a json of all the items belonging to the users which is given in the url
// the data returned is described by the item struct
func createItem(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var newItem SimpleItem
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

	json, err := json.MarshalIndent(newItem, "", "    ")

	// if conversion to json failed
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}
	_, err = db.Exec("insert into Inventory(UserID, TypeID) values (?, ?);", newItem.UserID, newItem.ItemType)

	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("add user error: %s", err)
		fmt.Fprintln(*w, err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}
func listUserItems(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	// row, err := db.Query("SELECT ItemID FROM Inventory WHERE UserID = ? ODER BY ItemID", r.PathValue("id"))
	row, err := db.Query(`
		SELECT Inventory.ItemID, Inventory.TypeID, ItemTypes.ItemName, ItemTypes.ItemDescription, ItemTypes.ImgURL, Marketplace.OfferID
		FROM Inventory
		INNER JOIN ItemTypes on Inventory.TypeID = ItemTypes.TypeID
        LEFT JOIN Marketplace on Inventory.ItemID = Marketplace.ItemID
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
		err := row.Scan(&item.ItemID, &item.TypeID, &item.ItemName, &item.ItemDescription, &item.ImgURL, &item.IsListed)

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
	//log.Printf(": %s", string(json))
	//log.Printf("test")
	// if conversion to json failed
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}
