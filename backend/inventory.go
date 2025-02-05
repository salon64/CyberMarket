package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "log"
	"net/http"
)

type Item struct{
	TypeID int
	ItemName string
	ItemDescription string
	ImgURL string
}

func listUserItems(w *http.ResponseWriter, r *http.Request, db *sql.DB){
	// row, err := db.Query("SELECT ItemID FROM Inventory WHERE UserID = ? ODER BY ItemID", r.PathValue("id"))
	row, err := db.Query("SELECT * FROM ItemTypes WHERE TypeID OPERATOR (SELECT ItemID FROM Inventory WHERE UserID = ? ODER BY ItemID)", r.PathValue("id"))
	if isErrLog(w, err) {
		return
	}

	defer row.Close()

	var items []Item
	for row.Next() {
		var item Item

		// write error and exit if scan fails
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