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

type TransactionInformation struct {
	TransID int
	Price int
	Date string
	ItemID int
	Buyer int
	Seller int
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

func displayTransactionslog(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var SQLStatement string 
	var row *sql.Rows
	var err error


	if r.PathValue("id") == "" {
		SQLStatement = `SELECT * FROM TransactionLog`
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