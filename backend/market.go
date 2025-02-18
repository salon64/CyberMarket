package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	// "log"
	"net/http"
)

type MarketplaceItems struct {
	ItemID int
	Price        int
}

type MarketplaceItemsInformation struct {
	ItemID int
	TypeID int
	UserID int

	ItemName        string
	ItemDescription *string
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
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)


	// TODO check if i own the item
	// TODO ERROR handling
	t, _ := db.Begin()
	res, err := t.Exec("insert into Marketplace(ItemID, Price, CreationDate) values (?, ?, now());", data.ItemID, data.Price)

	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("add listing error: %s", err)
		fmt.Fprintf(*w, "add listing error: %s", err.Error())
		return
	}
	// TODO: WE NEED A TRANSACTION HERE
	OfferID, err := res.LastInsertId()
	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(*w, err.Error())
		return
	}

	// return the id
	t.Commit()
	fmt.Fprintf(*w, "%d", OfferID)
}

func removeListingFromMarketplace(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var marketplaceItem MarketplaceItems

	if r.Body == nil {
		log.Print("body was nil")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&marketplaceItem)

	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "Error parsing json: %s", err.Error())
		log.Printf("Error parsing json: %s", err.Error())
		return
	}

	_, err = db.Exec("DELETE FROM Marketplace WHERE ItemID = ?;", r.PathValue("ItemID"))

	// if error write error and exit
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("error deleting marketplace listing: %s", err)
		fmt.Fprintln(*w, err.Error())
		return
	}
	fmt.Fprintln(*w, "removed listing")
}

func listMarketplaceItems(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data displayConstraints
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	log.Printf("%+v", data)
	if err != nil {
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)

	var orderBY string
	var SQLstatement string
	//mabe do some sprintf here to avoid this huge if else
	if data.SortBy == "Newest" {
		orderBY = "mp.OfferID"
		SQLstatement = `
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID
		order by mp.OfferID;
		`
	} else if data.SortBy == "Oldest" {
		orderBY = "mp.OfferID"
		SQLstatement = `
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID
		order by mp.OfferID DESC;
		`
	} else if data.SortBy == "Price_Ascending" {
		orderBY = "mp.Price"
		SQLstatement = `
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID
		order by mp.Price;
		`
	} else if data.SortBy == "Price_Descending" {
		orderBY = "mp.Price DESC"
		SQLstatement = `
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID
		order by mp.Price DESC;
		`
	} else if data.SortBy == "Alphabetically_Ascending" {
		orderBY = "it.ItemName"
		SQLstatement = `
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID
		order by it.ItemName;
		`
	} else if data.SortBy == "Alphabetically_Descending" {
		orderBY = "it.ItemName DESC"
		SQLstatement = `
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID
		order by it.ItemName DESC;
		`
	} else {
		// TODO ERROR HERE
		// OR DEFAULT HERE
	}
	log.Print(orderBY) // this is here since orderby needs to be used

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

	if isErrLog(w, err) {
		return
	}

	defer row.Close()

	var listings []MarketplaceItemsInformation
	for row.Next() {
		var listing MarketplaceItemsInformation
		// SELECT inv.ItemID, inv.TypeID, inv.UserID, u.Username, it.ItemName, it.ItemDescription, it.ImgURL, mp.OfferID, mp.Price, mp.CreationDate
		err := row.Scan(&listing.ItemID, &listing.TypeID, &listing.UserID, &listing.Username, &listing.ItemName, &listing.ItemDescription, &listing.ImgURL, &listing.OfferID, &listing.Price, &listing.CreationDate)
		if err != nil {
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(*w, err.Error())
			return
		}

		listings = append(listings, listing)
	}

	json, err := json.MarshalIndent(listings, "", "    ")

	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}

	// send json
	fmt.Fprint(*w, string(json))
}

func buyItem(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data MarketplaceItemsInformation
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)

	//TODO: check stuff

}