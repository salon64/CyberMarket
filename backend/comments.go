package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type AddComment struct {
	Grade   int
	UserID  *int
	Comment string
}

// Adds a comment to the type specified in the path
// if the user haven't bought the item before, an error would be returned with bad request as status code.
// The created comment will have a date of when the request was received.
func addComment(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var commentStruct AddComment

	// get the type id from the path
	itemTypeID, err := strconv.Atoi(r.PathValue("ItemTypeID"))
	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "Cant convert to ItemID ", err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&commentStruct)

	// if error on decoding
	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "Error decoding json data ", err.Error())
		return
	}

	// error on wrong formatted felids
	if commentStruct.Comment == "" || commentStruct.UserID == nil {
		sendAndLogError(w, http.StatusBadRequest, "Fields UserID or Comment dont exist or are an empty string")
		return
	}

	// no transaction is needed since adding a comment isn't critical
	// if fore some reason the user writes a review in the same times as it buys the item
	// an error could occur since the order of adding and checking is not guarantied
	// to save performance and avoid deadlocks no transaction is used

	// see if user have bought this item before,
	// since adding an item is the same as buying from null it will always detect if the user has the item
	row := db.QueryRow(`SELECT COUNT(1)
				FROM main_db.TransactionLog
				LEFT JOIN main_db.Inventory
				ON main_db.TransactionLog.ItemID = main_db.Inventory.ItemID
				WHERE
				main_db.TransactionLog.Buyer = ? AND main_db.Inventory.TypeID = ?;`,
		commentStruct.UserID, itemTypeID)

	// read if the row exist
	var exist int
	err = row.Scan(&exist)
	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "I should not happen ", err.Error())
		return
	}

	// if the the user haven't bought the item
	if exist == 0 {
		sendAndLogError(w, http.StatusBadRequest, "User haven't bought the item that it wants to review")
		return
	}

	// insert the comment
	db.Exec(`INSERT INTO main_db.TypeComments (UserID, TypeID, Grade, Comment, CreatedOn)
			VALUES (?,?,?,?, NOW())`,
		commentStruct.UserID, itemTypeID, commentStruct.Grade, commentStruct.Comment)

	fmt.Fprint(*w, "Added the comment")
}
