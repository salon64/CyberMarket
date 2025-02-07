package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

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

	defer row.Close()

	var items []Item
	for row.Next() {
		var item Item

		// write error and exit if scan fails
		err := row.Scan(&item.ItemID, &item.TypeID, &item.ItemName, &item.ItemDescription, &item.ImgURL)
		if err != nil {
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(*w, err.Error())
			return
		}

		items = append(items, item)
	}

	json, err := json.MarshalIndent(items, "", "    ")

	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}
