package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

	// is the user in the struct the same as the one posting or admin? 
	ok, _ := AuthByHeader(r,*commentStruct.UserID,db)
	if !ok {
		sendAndLogError(w,http.StatusForbidden,"auth failed")
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

type itemTypeReturn struct {
	Name      string
	ImgURL    *string
	ShortDesc *string
	DescURL   *string
	Comments  []pubComment
}

type pubComment struct {
	CommentID int
	UserName  string
	UserID    int
	Grade     int
	Comment   string
	PostedOn  string
}

// dont know if its the best idea to have it here, but its comment related ish
func getItemTypeInfo(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// convert the passed typeid to int
	id_str := r.PathValue("ItemTypeID")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		sendAndLogError(&w, http.StatusBadRequest, "Can't convert to a valid type id")
		return
	}
	log.Print("id ", id)

	// Get information about the item type
	// no transactions since its only reads and item types cant be deleted
	row := db.QueryRow(`SELECT
						main_db.ItemTypes.ItemName,
						main_db.ItemTypes.ImgURL,
						main_db.ItemTypes.ShortDescription,
						main_db.ItemTypes.ItemDescription
						FROM
						main_db.ItemTypes
						WHERE main_db.ItemTypes.TypeID = ?;`,
		id)

	// create the return struct
	ret := itemTypeReturn{}

	// scan the info into the ret struct
	err = row.Scan(&ret.Name, &ret.ImgURL, &ret.ShortDesc, &ret.DescURL)
	// if it cant find rows return that the itemid dont exist
	if errors.Is(err, sql.ErrNoRows) {
		sendAndLogError(&w, http.StatusNotFound, "Cant find the specified item")
		return
	}
	// if other error write error and return
	if err != nil {
		sendAndLogError(&w, http.StatusInternalServerError, "error getting item info", err.Error())
		return
	}

	log.Print(ret)

	rows, err := db.Query(`SELECT
							main_db.TypeComments.CommentID,
							main_db.Users.Username,
							main_db.Users.UserID,
							main_db.TypeComments.Grade,
							main_db.TypeComments.Comment,
							main_db.TypeComments.CreatedOn
							FROM main_db.TypeComments
							LEFT JOIN main_db.Users on main_db.TypeComments.UserID = main_db.Users.UserID
							WHERE main_db.TypeComments.TypeID = ?;`,
		id)
	if err != nil {
		sendAndLogError(&w, http.StatusInternalServerError, "error on getting items comment: ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var comment pubComment
		err := rows.Scan(&comment.CommentID, &comment.UserName, &comment.UserID, &comment.Grade, &comment.Comment, &comment.PostedOn)
		if err != nil {
			sendAndLogError(&w, http.StatusInternalServerError, "error while scanning comments: ", err.Error())
			return
		}

		ret.Comments = append(ret.Comments, comment)
	}

	json, err := json.MarshalIndent(ret, "", "    ")

	if err != nil {
		sendAndLogError(&w, http.StatusInternalServerError, err.Error())
		return
	}

	// send json
	fmt.Fprint(w, string(json))
}

// to delete a comment
func deleteComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	commentID, err := strconv.Atoi(r.PathValue("CommentID"))
	if err != nil {
		sendAndLogError(&w, http.StatusBadRequest, "can't convert ", r.PathValue("CommentID"), " to valid int")
		return
	}

	t, err := db.Begin()
	if err != nil {
		sendAndLogError(&w, http.StatusInternalServerError, "can't begin transaction: ", err.Error())
	}

	row := t.QueryRow("select UserID from main_db.TypeComments where CommentID = ?;",commentID)
	
	var ownerID int
	err = row.Scan(&ownerID)
	
	if errors.Is(err,sql.ErrNoRows) {
		t.Rollback()
		sendAndLogError(&w,http.StatusNotFound, "Cant find comment: ", err.Error())
		return
	} else if err != nil {
		t.Rollback()
		sendAndLogError(&w, http.StatusInternalServerError, "scan error: ", err.Error())
		return
	}

	ok, _ := AuthByHeader(r,ownerID,t)
	if !ok {
		t.Rollback()
		sendAndLogError(&w,http.StatusForbidden,"auth failed")
		return
	}

	_, err = t.Exec(`DELETE FROM main_db.TypeComments WHERE CommentID = ?;`, commentID)
	if err != nil {
		t.Rollback()
		sendAndLogError(&w,http.StatusInternalServerError, "error deleting comment: ", err.Error())
		return
	}

	t.Commit()
	fmt.Fprint(w, "Deleted comment")
}
