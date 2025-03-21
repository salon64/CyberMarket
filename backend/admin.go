package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	// "log"
	"net/http"
)

type ItemTypeInformation struct {
	ItemName         string
	ItemDescription  *string
	ImgURL           *string
	ShortDescription *string
}

type TransactionInformation struct {
	TransID int
	Price   int
	Date    string
	ItemID  int
	Buyer   int
	Seller  int
}

// creates a new item type
func createNewItemType(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data ItemTypeInformation
	AuthByHeader(r,-1,db)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "Error decoding json: ", err.Error())
		return
	}

	log.Printf("Adding a new item type %v", data)

	_, err = db.Exec(`INSERT INTO ItemTypes
						(ItemName, ItemDescription, ImgURL, ShortDescription)
						values (?, ?, ?, ?);`,
		data.ItemName, data.ItemDescription, data.ImgURL, data.ShortDescription)

	// if error write error and exit
	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "Error creating ItemType: ", err.Error())
		return
	}

	fmt.Fprint(*w, "success on creating item type")
}

func displayTransactionLogs(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	ownerID, err := strconv.Atoi(r.PathValue("id")) 
	
	if err != nil {
		ownerID = -1
	}

	// is user admin
	// no transaction needed since userid is used in the query and not indirectly referenced, 
	// therefore can not change before important our evaluation
	ok, _ := AuthByHeader(r,ownerID,db)
	if !ok {
		sendAndLogError(w,http.StatusForbidden,"auth failed")
		return
	}
	

	var SQLStatement string
	var row *sql.Rows

	if r.PathValue("id") == "all" {
		SQLStatement = `
		SELECT * FROM TransactionLog
		ORDER BY TransactionLog.Date DESC;`

		row, err = db.Query(SQLStatement)
	} else {
		SQLStatement = `
		SELECT * FROM TransactionLog
		WHERE TransactionLog.Buyer = ? OR TransactionLog.Seller = ?
		ORDER BY TransactionLog.Date DESC;
		`
		row, err = db.Query(SQLStatement, r.PathValue("id"), r.PathValue("id"))
	}
	// row, err := db.Query(SQLStatement, r.PathValue("id"), r.PathValue("id"))

	if isErrLog(w, err) {
		return
	}
	defer row.Close()

	var transactions []TransactionInformation
	for row.Next() {
		var transaction TransactionInformation
		err := row.Scan(&transaction.TransID, &transaction.Price, &transaction.Date, &transaction.ItemID, &transaction.Buyer, &transaction.Seller)
		if err != nil {
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(*w, err.Error())
			return
		}
		transactions = append(transactions, transaction)
	}

	json, err := json.MarshalIndent(transactions, "", "    ")

	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}
	// send json
	fmt.Fprint(*w, string(json))
}

type UserMoney struct {
	UserID int
	Money  int
}

func addMoneyToUser(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	// is user admin
	// no transaction needed since userid is used in the query and not indirectly referenced, 
	// therefore can not change before important our evaluation
	ok, _ := AuthByHeader(r,-1,db)
	if !ok {
		sendAndLogError(w,http.StatusForbidden,"auth failed")
		return
	}


	var data UserMoney
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "Error parsing json: ", err.Error())
		return
	}

	_, err = db.Exec("UPDATE Users SET Wallet = Wallet + ? WHERE UserID = ?", data.Money, data.UserID)
	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "error updating user wallet: ", err.Error())
		return
	}
	(*w).WriteHeader(http.StatusOK)
}

type SimpleItem struct {
	UserID   int
	ItemType int
}

// this functions writes out a json of all the items belonging to the users which is given in the url
// the data returned is described by the item struct
func createItem(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	// is user admin
	ok, _ := AuthByHeader(r,-1,db)
	if !ok {
		sendAndLogError(w,http.StatusForbidden, "auth failed")
	}

	var newItem SimpleItem
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newItem)

	if err != nil {
		sendAndLogError(w,http.StatusBadRequest, "Error parsing json: ", err.Error())
		return
	}

	_, err = db.Exec("insert into Inventory(UserID, TypeID) values (?, ?);", newItem.UserID, newItem.ItemType)

	if err != nil {
		sendAndLogError(w,http.StatusInternalServerError, "adding ItemType returned error: ", err.Error())
		return
	}

	// send json
	(*w).WriteHeader(http.StatusOK)
}
