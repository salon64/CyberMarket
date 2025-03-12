package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type inCheckout struct {
	InCheckout int
}

// returns all items in users cart. same format as list marketplace
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
		err := row.Scan(
			&listing.ItemID,
			&listing.TypeID,
			&listing.UserID,
			&listing.Username,
			&listing.ItemName,
			&listing.ShortDescription,
			&listing.ImgURL,
			&listing.OfferID,
			&listing.Price,
			&listing.CreationDate)
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
	var data UserStruct

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

// this function returns if the if the item is in the users cart. zero for false 1 for true
func checkCart(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data UserStruct

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
	fmt.Fprint(*w, string(json))
}

// Removes an item forom the users cart, pahtvalue is offerid and struct contains userod
func removeFromCart(w *http.ResponseWriter, r *http.Request, db *sql.DB) {
	var data UserStruct

	if r.Body == nil {
		log.Print("body was nil")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendAndLogError(w, http.StatusBadRequest, "error decoding: ", err.Error())
		return
	}

	_, err = db.Exec("DELETE FROM ShoppingCart WHERE (UserID, OfferID) = (?, ?);", data.UserID, r.PathValue("OfferID"))

	// if error write error and exit
	if err != nil {
		sendAndLogError(w, http.StatusInternalServerError, "shopping basket insertion error: ", err.Error())
		return
	}

	fmt.Fprintf(*w, "%s was deleted from cart", r.PathValue("OfferID"))
}

type cartStruckt struct {
	price  int
	owner  int
	itemid int
}

func buyCart(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// begin a sql transaction
	t, err := db.Begin()

	// TODO use other util func here
	if isErrLog(&w, err) {
		log.Print("HERE")

		return
	}

	// get price of each item in cart

	var cartPrice int
	err = t.QueryRow(`SELECT SUM(Marketplace.Price)
					FROM ShoppingCart
					INNER JOIN Marketplace ON ShoppingCart.OfferID = Marketplace.OfferID
					WHERE ShoppingCart.UserID = ?;`,
		r.PathValue("UserID"),
	).Scan(&cartPrice)
	if isErrLog(&w, err) {
		log.Print("HERE 2")

		return
	}

	// get user wallet
	var user_wallet int
	err = t.QueryRow("SELECT Users.Wallet FROM Users WHERE UserID = ?;", r.PathValue("UserID")).Scan(&user_wallet)
	if isErrLog(&w, err) {
		log.Print("HERE 3")

		return
	}

	// check for sufficient funds
	if user_wallet < cartPrice {
		log.Print("HERE 4")

		t.Rollback()
		fmt.Fprint(w, "Insufficient funds to buy items in cart")
		return
	}

	// update user wallet
	_, err = t.Exec(`UPDATE Users SET Wallet = Wallet - ? WHERE Users.UserID = ? ;`, cartPrice, r.PathValue("UserID"))
	if err != nil {
		log.Print("HERE 5")

		t.Rollback()
		return
	}

	// for each item in shopping cart, get the owner, price and itemid
	rows, err := t.Query(`SELECT Marketplace.Price, Inventory.ItemID, Inventory.UserID
		FROM ShoppingCart
		LEFT JOIN Marketplace ON ShoppingCart.OfferID = Marketplace.OfferID
		LEFT JOIN Inventory ON Marketplace.ItemID = Inventory.ItemID
		WHERE ShoppingCart.UserID = ?;`, r.PathValue("UserID"))
	if isErrLog(&w, err) {
		log.Print("HERE 6")

		t.Rollback()
		return
	}
	defer rows.Close()

	// for each item in users cart

	cart := make([]cartStruckt, 0, 10)
	for rows.Next() {
		var price, owner, itemID int
		err = rows.Scan(&price, &itemID, &owner)
		if isErrLog(&w, err) {
			log.Print("HERE 7")

			t.Rollback()
			return
		}
		cart = append(cart, cartStruckt{price: price, owner: owner, itemid: itemID})
	}

	for _, cart_item := range cart {

		// update owner
		_, err = t.Exec(`UPDATE Inventory SET UserID = ? WHERE Inventory.ItemID = ? ;`, r.PathValue("UserID"), cart_item.itemid)
		if err != nil {
			log.Print("HERE 8")

			t.Rollback()
			return
		}

		// add funds to item owner
		_, err = t.Exec(`UPDATE Users SET Wallet = Wallet + ? WHERE Users.UserID = ? ;`, cart_item.price, cart_item.owner)
		if err != nil {
			log.Print("HERE 9")

			t.Rollback()
			return
		}

		// delete listing
		// will also delete from the cart as the relation is set to cascade
		_, err = db.Exec("DELETE FROM Marketplace WHERE ItemID = ?;", cart_item.itemid)
		if err != nil {
			log.Print("HERE 10")

			t.Rollback()
			return
		}
		// create transaction log
		_, err = t.Exec(`INSERT INTO TransactionLog (Price, Date, ItemID, Buyer, Seller) 
			VALUES (?, NOW(), ?, ?, ?);`, cart_item.price, cart_item.itemid, r.PathValue("UserID"), cart_item.owner)
		if err != nil {
			log.Printf("HERE 11 %s", err.Error())

			t.Rollback()
			return
		}
	}

	err = t.Commit()
	if isErrLog(&w, err) {
		log.Print("HERE 12")

		return
	}
	fmt.Fprint(w, "bought items from cart")

}
