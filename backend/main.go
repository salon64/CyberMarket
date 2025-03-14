package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// if error is not nil, then the error is written to log, response and true is returned
func isErrLog(w *http.ResponseWriter, err error) bool {
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, err.Error())
		log.Print(err.Error())
		return true
	}
	return false
}

func main() {
	// set upp logger
	log.SetFlags(log.Flags() | log.Lshortfile)

	cfg := mysql.Config{
		User:      os.Getenv("DBUSER"),
		Passwd:    os.Getenv("DBPASS"),
		Net:       "tcp",
		ParseTime: true,
		Addr:      os.Getenv("DBHOST"),
		DBName:    "main_db", // this is the name that is defined
	}
	var err error // used here so that it makes db global

	// try to connect 30 times each waiting 2 seconds
	log.Printf("Connecting to database at host %s, with DSN %s", os.Getenv("DBHOST"), cfg.FormatDSN())
	for i := 1; i < 31; i++ {

		// try to connect with cfg as connection options
		db, err = sql.Open("mysql", cfg.FormatDSN())

		// if connection attempt returned error log
		if err != nil {
			log.Printf("Error connecting, attempt %d/30. %s. ", i, err.Error())
			time.Sleep(time.Second * 2)
		} else {
			log.Print("Connection established")
			break
		}
	}

	// if error still exist after 30 connection tries, abort
	if err != nil {
		log.Fatal("aborting failed to establish connection")
	}
	defer db.Close() // close db when main exits

	// try to ping 30 times each waiting 2 seconds

	var pingErr error // test
	for i := 1; i < 31; i++ {
		// pinging db
		pingErr = db.Ping()

		// if ping
		if pingErr != nil {
			log.Printf("Error pinging, attempt %d/30. %s. ", i, pingErr.Error())
			time.Sleep(time.Second * 2)
		} else {
			log.Print("Connection verified")
			break
		}
	}
	// if error still exist after 30 ping attempts, abort
	if pingErr != nil {
		log.Fatal("aborting failed to ping")
	}

	// TODO set max life time and other


	http.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-control-allow-methods", "POST, GET")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	})
	
	// function that returns an array containing PubUser struct {id,name} in json format
	http.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		// allow CORS
		enableCors(&w)
		listAllUsers(&w, r, db)
	})
	
	http.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		addUser(&w, r, db)
	})
	
	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		userLogin(&w, r, db)
	})
	
	http.HandleFunc("POST /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		log.Print("I got post ", time.Now().GoString())
		enableCors(&w)
		updateUserInfo(&w, r, db)
	})

	http.HandleFunc("GET /user/getMoney/{id}", func(w http.ResponseWriter, r *http.Request) {
		enableCors((&w))
		getUserMoney(&w, r, db)
	})

	http.HandleFunc("POST /user/AddMoney", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		addMoneyToUser(&w, r, db)
	})

	http.HandleFunc("GET /inventory/{id}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		listUserItems(&w, r, db)
	})

	// adding a comment to an item type
	http.HandleFunc("POST /ItemType/{ItemTypeID}", func(w http.ResponseWriter, r *http.Request) {

		enableCors(&w)
		addComment(&w, r, db)
	})

	http.HandleFunc("GET /ItemType/{ItemTypeID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		getItemTypeInfo(w, r, db)
	})

	//TODO: update to DELETE, and handle option
	http.HandleFunc("GET /comment/deletecomment/{CommentID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		deleteComment(w, r, db)
	})

	http.HandleFunc("POST /Marketplace/displayMarket", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		listMarketplaceItems(&w, r, db)
	})

	http.HandleFunc("GET /displayTransactionslog/{id}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		displayTransactionLogs(&w, r, db)
	})

	//TODO change to OfferID
	http.HandleFunc("POST /Marketplace/buy/{ItemID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		err := buyItem(&w, r, db)
		if err != nil {
			log.Print(err.Error())
			fmt.Fprint(w, err.Error())
		}
	})

	http.HandleFunc("POST /Marketplace/addListing", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		addListingToMarketplace(&w, r, db)
	})

	http.HandleFunc("POST /Admin/CreateNewItemType", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		createNewItemType(&w, r, db)
	})
	http.HandleFunc("GET /Marketplace/displayCart/{UserID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		displayCart(&w, r, db)
	})
	http.HandleFunc("POST /Marketplace/checkCart/{OfferID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		checkCart(&w, r, db)
	})
	http.HandleFunc("POST /Marketplace/addToCart/{OfferID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		addToCart(&w, r, db)
	})
	http.HandleFunc("POST /Marketplace/removeFromCart/{OfferID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		removeFromCart(&w, r, db)
	})
	http.HandleFunc("GET /Marketplace/cartCheckout/{UserID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		buyCart(w, r, db)
	})
	http.HandleFunc("GET /Marketplace/removeListing/{ItemID}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		removeListingFromMarketplace(&w, r, db)
	})

	http.HandleFunc("POST /Marketplace/CreateItem", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		createItem(&w, r, db)
	})

	http.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		row, err := db.Query("SHOW TABLES;")

		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			for row.Next() {
				var table string
				err := row.Scan(&table)
				if err != nil {
					fmt.Fprintln(w, err.Error())
				} else {
					fmt.Fprintln(w, table)
				}
			}
		}
	})

	// err = http.ListenAndServeTLS(":5687",os.Getenv("CERT_FILE_PATH"),os.Getenv("CERT_KEY_FILE_PATH"),nil)
	// if err != nil {
	// 	log.Print("HTTPS ", err.Error())
	// }

	err = http.ListenAndServe(":5687", nil)
	if err != nil {
		log.Print("HTTP ", err.Error())
	}
}
