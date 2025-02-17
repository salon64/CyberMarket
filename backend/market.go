package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"encoding/json"
	// "log"
	"net/http"
)

type  MarketplaceItems struct{
	ItemID: int;
	// OfferID: int;
	Price: int;
	CreationDate: string;
}

type MarketplaceItemsInformation struct{
	ItemID: int;
	TypeID: int;
	UserID: int;

	ItemName: string;
	ItemDescription: string;
	ImgURL: string;
	
	OfferID: int;
	Price: int;
	CreationDate: string;

	Username: string;
}





func addListingToMarketplace(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var marketplaceItem MarketplaceItems

	if r.Body == nil {
		log.Print("body was nil")
		return
	}


	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&useradd)

	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "Error parsing json: %s", err.Error())
		log.Printf("Error parsing json: %s", err.Error())
		return
	}
	// TODO check if i own the item 
	t := db.Begin()
	res, err := t.Exec("insert into Marketplace(ItemID, Price, CreationDate) values (?, ?, ?);", marketplaceItem.ItemID, marketplaceItem.Price, marketplaceItem.CreationDate)

	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("add listing error: %s", err)
		fmt.Fprintf(*w,"add listing error: %s", err.Error())
		return
	}
	// TODO WE NEED A TRANSACTION HERE
	OfferID, err := res.LastInsertId()
	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(*w, err.Error())
		return
	}

	// return the id
	t.commit()
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
	err := decoder.Decode(&useradd)

	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "Error parsing json: %s", err.Error())
		log.Printf("Error parsing json: %s", err.Error())
		return
	}

	res, err := Exec("DELETE FROM Marketplace WHERE OfferID = ?;", 1)

	// if error write error and exit
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("error deleting marketplace listing: %s", err)
		fmt.Fprintln(*w, err.Error())
		return
	}
	fmt.Fprintln(w,"removed listing")
}


func listMarketplaceItems(w *http.ResponseWriter, _ *http.Request, db *sql.DB) {
	var data MarketplaceItemsInformation{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)

	var orderBY 
	if data

	

	var SQLstatement := "
		SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID
		order by  ?;
	"

	
	row, err := db.Query(SQLstatement)

	if isErrLog(w, err) {
		return
	}

	defer row.Close()

	var listings []MarketplaceItemsInformation
	for row.Next() {
		var listing MarketplaceItemsInformation

		err := row.Scan(&listing.ItemID, &listing.TypeID, &listing.ItemName, &listing.ItemDescription, &listing.ImgURL, &listing.OfferID, &listing.Price, &listing.CreationDate)
		if err != nil {
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(*w, err.Error())
			return
		}

		listings = append(listings, listing)
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

func buyItem(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data MarketplaceItemsInformation{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("error decoding: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error decoding: %s", err.Error())
		return
	}
	log.Printf("with data %v", data)




}

func sellItem() {
	
}

