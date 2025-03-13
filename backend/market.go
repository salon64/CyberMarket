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

type UserStruct struct {
	UserID int
}
type MarketplaceItems struct {
	ItemID int
	Price  int
}

type MarketplaceItemsInformation struct {
	ItemID int
	TypeID int
	UserID int

	ItemName        string
	ShortDescription *string
	ImgURL          *string

	OfferID      int
	Price        int
	CreationDate string

	Username string
}
type displayConstraints struct {
	SortBy string
	Search string
}

func addListingToMarketplace(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data MarketplaceItems

	if r.Body == nil {
		log.Print("body was nil")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "Error parsing json: ", err.Error())
		return
	}

	// start an transaction
	t, err := db.Begin()
	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "Error starting transaction: ", err.Error())
		return
	}

	// test if user is allowed to place listing
	row := t.QueryRow("select UserID from main_db.Inventory where ItemID = ?;",data.ItemID)
	
	var ownerID int
	err = row.Scan(&ownerID)
	
	if errors.Is(err,sql.ErrNoRows) {
		t.Rollback()
		sendAndLogError(w,http.StatusNotFound, "Cant find item: ", err.Error())
		return
	} else if err != nil {
		t.Rollback()
		sendAndLogError(w, http.StatusInternalServerError, "scan error: ", err.Error())
		return
	}

	// this needs to be part of the transaction, 
	// imagine if not
	// 1. Alice buys an item from Bob, a transaction is started to Transfer owner ship
	// 2. Bob add begins adding a listing for the same item Alice is buying
	// 3. Bobs auth check is passed when adding
	// 4. Alice transaction is completed
	// 5. Bob have passed auth and now an listing is added for the item now Alice owns
	// 6. Bob chose the price of zero euro dollar and buys the item again from Alice
	//
	// The result of this is that Alice Paid for an item that was immediately "stolen"
	ok, _ := AuthByHeader(r,ownerID,t)
	if !ok {
		t.Rollback()
		sendAndLogError(w,http.StatusForbidden,"auth failed")
		return
	}

	// add marketplace listing
	res, err := t.Exec("insert into Marketplace(ItemID, Price, CreationDate) values (?, ?, now());", data.ItemID, data.Price)

	// if error write error and exit
	if err != nil {
		t.Rollback()
		sendAndLogError(w, http.StatusInternalServerError, "Error inserting listing: ", err.Error())
		return
	}

	OfferID, err := res.LastInsertId()
	// if error write error and exit
	if err != nil {
		t.Rollback()
		sendAndLogError(w, http.StatusInternalServerError, "Error getting transaction id: ", err.Error())
		return
	}

	// return the id

	err = t.Commit()
	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "Error committing listing: ", err.Error())
	}

	fmt.Fprintf(*w, "added listing with id: %d", OfferID)
}

func removeListingFromMarketplace(w *http.ResponseWriter, r *http.Request, db *sql.DB) {


	// start an transaction
	t, err := db.Begin()
	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "Error starting transaction: ", err.Error())
		return
	}

	// test if user is allowed to delete listing
	// get the listing owner
	row := t.QueryRow("select UserID from main_db.Inventory where ItemID = ?;",r.PathValue("ItemID"))
	
	var ownerID int
	err = row.Scan(&ownerID)
	
	if errors.Is(err,sql.ErrNoRows) {
		t.Rollback()
		sendAndLogError(w,http.StatusNotFound, "Cant find item: ", err.Error())
		return
	} else if err != nil {
		t.Rollback()
		sendAndLogError(w, http.StatusInternalServerError, "scan error: ", err.Error())
		return
	}
	// is the user the owner or an admin
	ok, _ := AuthByHeader(r,ownerID,t)
	if !ok {
		t.Rollback()
		sendAndLogError(w,http.StatusForbidden,"auth failed")
		return
	}

	// delete the listing
	_, err = t.Exec("DELETE FROM Marketplace WHERE ItemID = ?;", r.PathValue("ItemID"))

	// if error write error and exit
	if err != nil {
		t.Rollback()
		sendAndLogError(w, http.StatusInternalServerError, "Error deleting listing: ", err.Error())
		return
	}

	// commit the deletion
	err = t.Commit()
	if err != nil {
		sendAndLogError(w,http.StatusInternalServerError, "error committing transaction: ", err.Error())
		return
	}

	fmt.Fprintln(*w, "removed listing")
}

func listMarketplaceItems(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data displayConstraints
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "Error parsing json: ", err.Error())
		return
	}

	var SQLstatement string
	SQLstatement = `
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ShortDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID `

	// maybe do some sprintf here to avoid this huge if else
	if data.SortBy == "Newest" {
		SQLstatement += "order by CreationDate;"

	} else if data.SortBy == "Oldest" {
		SQLstatement += `
		order by CreationDate DESC;
		`
	} else if data.SortBy == "Price_Ascending" {
		SQLstatement += `
		order by mp.Price;
		`
	} else if data.SortBy == "Price_Descending" {
		SQLstatement += `
		order by mp.Price DESC;
		`
	} else if data.SortBy == "Alphabetically_Ascending" {
		SQLstatement += `
		order by it.ItemName;
		`
	} else if data.SortBy == "Alphabetically_Descending" {
		SQLstatement += `
		order by it.ItemName DESC;
		`
	} else {
		SQLstatement += `;`
	}

	// var SQLstatement := `
	// 	SELECT inv.ItemID, inv.TypeID, inv.UserID,
	// 	u.Username,
	// 	it.ItemName, it.ItemDescription, it.ImgURL,
	// 	mp.OfferID, mp.Price, mp.CreationDate
	// 	FROM Marketplace mp
	// 	INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
	// 	INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
	// 	INNER JOIN Users u ON u.UserID = inv.UserID
	// 	order by  ?;
	// `

	row, err := db.Query(SQLstatement)

	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "Error when querying marketplace: ", err.Error())
		return
	}

	defer row.Close()

	var listings []MarketplaceItemsInformation
	for row.Next() {
		var listing MarketplaceItemsInformation
		// SELECT inv.ItemID, inv.TypeID, inv.UserID, u.Username, it.ItemName, it.ItemDescription, it.ImgURL, mp.OfferID, mp.Price, mp.CreationDate
		err := row.Scan(&listing.ItemID, &listing.TypeID, &listing.UserID, &listing.Username, &listing.ItemName, &listing.ShortDescription, &listing.ImgURL, &listing.OfferID, &listing.Price, &listing.CreationDate)
		if err != nil {
			sendAndLogError(w, http.StatusInternalServerError, "Error scanning: ", err.Error())
			return
		}

		listings = append(listings, listing)
	}

	json, err := json.MarshalIndent(listings, "", "    ")

	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "error encoding return: ", err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}

type BuyStruct struct {
	UserID int
}

func buyItem(w *http.ResponseWriter, r *http.Request, db *sql.DB) error {
	//TODO auth by token
	var data BuyStruct
	itemID, err := strconv.Atoi(r.PathValue("ItemID"))
	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		return err
	}

	if r.Body == nil {
		(*w).WriteHeader(http.StatusBadRequest)
		return errors.New("body was empty")
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)

	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		return err
	}

	// no transaction needed since userid is used in the query and not indirectly referenced, 
	// therefore can not change before important our evaluation
	ok, _ := AuthByHeader(r,data.UserID, db)
	if !ok {
		(*w).WriteHeader(http.StatusForbidden)
		return errors.New("auth failed")
	}


	// begin transaction
	t, err := db.Begin()
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}

	// get the price and seller
	row := t.QueryRow(`SELECT Marketplace.Price, Inventory.UserID
						FROM main_db.Marketplace
						LEFT JOIN Inventory on  Inventory.ItemID = Marketplace.ItemID
						Where Marketplace.ItemID = ? ;`, itemID)
	var price int
	var seller int
	err = row.Scan(&price, &seller)
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}

	row = t.QueryRow("SELECT Wallet FROM main_db.Users WHERE UserID = ?;", data.UserID)
	var wallet int
	err = row.Scan(&wallet)
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	// test if user have enough funds
	if wallet < price {
		t.Rollback()
		return errors.New("not enough funds")
	}
	// update users funds and item owner
	_, err = t.Exec(`UPDATE Users SET Wallet = Wallet - ? WHERE Users.UserID = ? ;`, price, data.UserID)
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = t.Exec(`UPDATE Inventory SET UserID = ? WHERE Inventory.ItemID = ? ;`, data.UserID, itemID)
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = t.Exec(`UPDATE Users SET Wallet = Wallet + ? WHERE Users.UserID = ? ;`, price, seller)
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = db.Exec("DELETE FROM Marketplace WHERE ItemID = ?;", itemID)
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	// create transaction log
	_, err = t.Exec(`INSERT INTO TransactionLog (Price, Date, ItemID, Buyer, Seller) 
        VALUES (?, NOW(), ?, ?, ?);`, price, itemID, data.UserID, seller)
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	err = t.Commit()
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}

	fmt.Fprint(*w, "Success")
	return nil
}
