package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	// "log"
	"net/http"
)
type UserStruc struct{
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
	ItemDescription *string
	ImgURL          *string

	OfferID      int
	Price        int
	CreationDate string

	Username string
}
type inCheckout struct {
	InCheckout int
}
type displayConstraints struct {
	SortBy string
	Search string
}
// --------------------- SHOPPING CART --------------------------------------------

func displayCart(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	row, err := db.Query((`SELECT inv.ItemID, inv.TypeID, inv.UserID,
		u.Username,
		it.ItemName, it.ItemDescription, it.ImgURL, 
		mp.OfferID, mp.Price, mp.CreationDate
		FROM Marketplace mp
		INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
		INNER JOIN ItemTypes it ON inv.TypeID = it.TypeID
		INNER JOIN Users u ON u.UserID = inv.UserID 
		WHERE mp.OfferID IN (SELECT OfferID FROM ShoppingCart WHERE UserID = ?);`), r.PathValue("UserID"))

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

// Recieves UserID and OfferID and adds the offer to the cart table in DB
func addToCart(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data UserStruc

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

	t, _ := db.Begin()
	_, err = t.Exec("insert into ShoppingCart(UserID, OfferID) values (?, ?);", data.UserID, r.PathValue("OfferID"))

	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("ShoppingCart insertion error: %s", err)
		fmt.Fprintf(*w, "ShoppingCart insertion error: %s", err.Error())
		return
	}
	t.Commit()

	fmt.Fprintf(*w, "%s", r.PathValue("OfferID"))

}
func checkCart(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data UserStruc

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

	t, _ := db.Begin()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM ShoppingCart WHERE UserID = ? AND OfferID = ?", data.UserID, r.PathValue("OfferID")).Scan(&count)

	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("ShoppingCart check query error: %s", err)
		fmt.Fprintf(*w, "ShoppingCart check query error: %s", err.Error())
		return
	}

	t.Commit()
	tmp := inCheckout{InCheckout: count}
	json, err := json.MarshalIndent(tmp, "", "    ")

	// write error and exit if json fails
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		return
	}

	// send json
	log.Printf("string test")
	log.Println(string(json))
	fmt.Fprint(*w, string(json))


	// return tmp value to original html fetch post request
}
func removeFromCart(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data UserStruc

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

	t, _ := db.Begin()
	_, err = t.Exec("DELETE FROM ShoppingCart WHERE (UserID, OfferID) = (?, ?);", data.UserID, r.PathValue("OfferID"))

	// if error write error and exit
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("ShoppingCart insertion error: %s", err)
		fmt.Fprintf(*w, "ShoppingCart insertion error: %s", err.Error())
		return
	}
	t.Commit()
	
	fmt.Fprintf(*w, "%s", r.PathValue("OfferID"))
}
func cartCheckout(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Returns illegal rows, if any
	rows, err := db.Query(`SELECT sc.OfferID FROM ShoppingCart sc 
	WHERE sc.UserID = ? AND NOT EXISTS 
	(SELECT 1 FROM Marketplace mp WHERE 
	mp.OfferID = sc.OfferID);`, r.PathValue("UserID"))

	if rows.Next() { // Item doesn't exist in marketplace anymore, delete item from cart
		log.Printf("Item doesn't exist anymore")
		fmt.Fprintf(*w, "Item doesn't exist anymore")
		_, err := db.Exec(`DELETE FROM ShoppingCart 
		WHERE UserID = ? 
		AND NOT EXISTS (
			SELECT 1 
			FROM Marketplace mp 
			WHERE mp.OfferID = ShoppingCart.OfferID
		);`, r.PathValue("UserID"))
		log.Printf("test1")
		if err != nil { // Shouldn't be able to get an error here
			log.Printf("impossible error (removing cart item): %s", err.Error())
			(*w).WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(*w, "impossible error (removing cart item): %s", err.Error())
			return
		}
		return
	}
	if err != nil {
		log.Printf("error querying row (cartCheckout): %s", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(*w, "error querying row (cartCheckout): %s", err.Error())
		return
	}
	defer rows.Close()
	rows, err = db.Query(`SELECT OfferID FROM ShoppingCart WHERE UserID = ?`, r.PathValue("UserID"))
	if err != nil {
		log.Printf("error reading offerID rows: %s", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(*w, "error reading offerID rows: %s", err.Error())
		return
	}
	var totCost int
	var myWallet int
	// Get totalcost
	tmp := db.QueryRow(`SELECT SUM(mp.Price) AS TotalCost
		FROM ShoppingCart sc
		INNER JOIN Marketplace mp ON sc.OfferID = mp.OfferID
		WHERE sc.UserID = ?;`, r.PathValue("UserID"))
	err = tmp.Scan(&totCost)

	if err != nil {
		log.Printf("error reading totCost: %s", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(*w, "error reading totCost: %s", err.Error())
		return
	}

	// Get user wallet
	tmp = db.QueryRow(`SELECT Wallet FROM Users WHERE UserID = ?`, r.PathValue("UserID"))
	err = tmp.Scan(&myWallet)
	if err != nil {
		log.Printf("error reading buyer wallet: %s", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(*w, "error reading buyer wallet: %s", err.Error())
		return
	}
	if myWallet < totCost {
		log.Printf("Not enough money in wallet")
		fmt.Fprintf(*w, "Not enough money in wallet")
		return
	}

	t, err := db.Begin()
	if err != nil {
		log.Printf("error couldn't begin transaction: %s", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(*w, "error couldn't begin transaction: %s", err.Error())
		return
	}
	for rows.Next() {
		var currOfferID int
		var currPrice int
		var seller int
		var currItemID int

		err := rows.Scan(&currOfferID)
		if err != nil {
			t.Rollback()
			log.Printf("error reading current offerID: %s", err.Error())
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(*w, "error reading current offerID: %s", err.Error())
			return
		}

		tmp := t.QueryRow(`SELECT inv.UserID AS SellerID
			FROM Marketplace mp
			INNER JOIN Inventory inv ON mp.ItemID = inv.ItemID
			WHERE mp.OfferID = ?;`, currOfferID)
		err = tmp.Scan(&seller)
		if err != nil {
			t.Rollback()
			log.Printf("error reading seller: %s", err.Error())
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(*w, "error reading seller: %s", err.Error())
			return
		}
		tmp = t.QueryRow(`SELECT Price FROM Marketplace WHERE OfferID = ?`, currOfferID)
		err = tmp.Scan(&currPrice)
		if err != nil {
			t.Rollback()
			log.Printf("error reading seller: %s", err.Error())
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(*w, "error reading seller: %s", err.Error())
			return
		}

		tmp = t.QueryRow(`SELECT ItemID FROM Marketplace WHERE OfferID = ?`, currOfferID)
		err = tmp.Scan(&currItemID)
		if err != nil {
			t.Rollback()
			log.Printf("error reading itemID: %s", err.Error())
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(*w, "error reading itemID: %s", err.Error())
			return
		}
		// Remove money from buyer wallet
		_, err = t.Exec("UPDATE Users SET Wallet = Wallet - ? WHERE UserID = ?", currPrice, r.PathValue("UserID"))
		if err != nil {
			t.Rollback()
			log.Printf("error updating buyer wallet: %s", err.Error())
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(*w, "error updating buyer wallet: %s", err.Error())
			return
		}
		// Add money to seller wallet
		_, err = t.Exec("UPDATE Users SET Wallet = Wallet + ? WHERE UserID = ?", currPrice, seller)
		if err != nil {
			t.Rollback()
			log.Printf("error updating seller wallet: %s", err.Error())
			(*w).WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(*w, "error updating seller wallet: %s", err.Error())
			return
		}
		_, err = t.Exec(`UPDATE Inventory SET UserID = ? WHERE Inventory.ItemID = ? ;`, r.PathValue("UserID"), currItemID)
		if err != nil {
			t.Rollback()
			(*w).WriteHeader(http.StatusInternalServerError)
			return
		}
		// create transaction log
		_, err = t.Exec(`INSERT INTO TransactionLog (Price, Date, ItemID, Buyer, Seller) 
			VALUES (?, NOW(), ?, ?, ?);`, currPrice, currItemID, r.PathValue("UserID"), seller)
		if err != nil {
			t.Rollback()
			(*w).WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("DELETE FROM Marketplace WHERE ItemID = ?;", currItemID)
		if err != nil {
			t.Rollback()
			(*w).WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	_, err = db.Exec("DELETE FROM ShoppingCart WHERE UserID = ?;", r.PathValue("UserID"))
	if err != nil {
		t.Rollback()
		(*w).WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Commit()
	log.Printf("Successfully checked out")
	fmt.Fprint(*w, "Success")
}
// --------------------------------------------------------------------------------
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
	_, err = db.Exec("UPDATE Inventory SET IsListed = 1 WHERE ItemID = ?;", data.ItemID)
	if err != nil {
		log.Printf("error updating IsListed to true: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error updating IsListed to true: %s", err.Error())
		return
	}
	fmt.Fprintf(*w, "%d", OfferID)
}

func removeListingFromMarketplace(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, err := db.Exec("DELETE FROM Marketplace WHERE ItemID = ?;", r.PathValue("ItemID"))
	log.Printf("test")
	log.Printf("%s", r.PathValue("ItemID"))
	log.Printf("test")
	// if error write error and exit
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("error deleting marketplace listing: %s", err)
		fmt.Fprintln(*w, err.Error())
		return
	}

	_, err = db.Exec("UPDATE Inventory SET IsListed = 0 WHERE ItemID = ?;", r.PathValue("ItemID"))
	log.Print(r.PathValue("ItemID"))
	if err != nil {
		log.Printf("error updating IsListed to false: %s", err.Error())
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(*w, "error updating IsListed to false: %s", err.Error())
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
	log.Printf("%+v", data)
	log.Print(itemID)
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
