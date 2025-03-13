package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func listUserItems(w *http.ResponseWriter, r *http.Request, db *sql.DB) {

	ownerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "Cant convert to ItemID ", err.Error())
		return
	}

	// no need for transaction
	// since the owner id
	ok, _ := AuthByHeader(r,ownerID,db)
	if !ok {
		sendAndLogError(w,http.StatusForbidden,"auth failed")
		return
	}
	

	row, err := db.Query(`
		SELECT Inventory.ItemID, Inventory.TypeID, ItemTypes.ItemName, ItemTypes.ShortDescription, ItemTypes.ImgURL, Marketplace.OfferID
		FROM Inventory
		INNER JOIN ItemTypes on Inventory.TypeID = ItemTypes.TypeID
        LEFT JOIN Marketplace on Inventory.ItemID = Marketplace.ItemID
		where Inventory.UserID = ?;`,
		ownerID)

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
			sendAndLogError(w,http.StatusInternalServerError, "error scanning user items: ", err.Error())
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
		sendAndLogError(w, http.StatusInternalServerError, "error encoding return: ", err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}
