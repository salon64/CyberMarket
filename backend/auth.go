package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

func AuthByHeader(r *http.Request, resourceOwner int, db Queryer) (bool, error) {
	auth_row := r.Header.Get("Authorization")

	token, found := strings.CutPrefix(auth_row, "Bearer ")
	if !found {
		log.Print("auth is not of type Bearer")
		return false, errors.New("auth is not of type Bearer")
	}

	if token == "" {
		log.Print("token was empty")
		return false, errors.New("token was empty")
	}

	return AuthByToken(token, resourceOwner, db)
}

// This returns true if the token passes is valid and the token owner matches the resourceOwner
// Note currently that if the suer is an admin, access is always given
func AuthByToken(token string, resourceOwner int, db Queryer) (bool, error) {
	// select user ID user role and creation data of token
	row := db.QueryRow(`
		SELECT TokenTable.UserID, Users.Role, TokenTable.CreatedOn
		FROM TokenTable
		INNER JOIN Users on Users.UserID = TokenTable.UserID
		WHERE TokenTable.Token = UUID_TO_BIN(?);`,
		token,
	)

	var userID int
	var role int
	var time time.Time
	
	// read the values
	err := row.Scan(&userID, &role, &time)
	
	// the possible errors are normal scan and row missing caused by there not being shush a token
	if err != nil {
		log.Print("log error: ", err.Error())
		return false, err
	}
	// if user has a admin role
	if role != 0 {
		log.Println("user is admin, auth complete")
		return true, nil
		// if the user own the requested resource
	} else if userID == resourceOwner {
		log.Println("user ", userID, " is the owner of ", resourceOwner, ", auth complete")
		return true, nil
		// if not return false
	} else {
		log.Println("user ", userID, " is not the owner of ", resourceOwner, ", auth failed")
		return false, nil
	}

}

// This struct contains the return of CreateToken
type LoginReturn struct {
	UserID int
	Role   int
	Token  string
}

// If the username and password combination exist a string of the uuid userid is returned
func CreateToken(userName string, password string, db *sql.DB) (LoginReturn, error) {
	// TODO error handling
	// begin the transaction
	t, _ := db.Begin()
	// see if users exist
	row := t.QueryRow("SELECT userID, Role FROM Users WHERE Username = ? AND `Password(Hash)` = ?;", userName, password)

	// read user id
	var userID int
	var role int
	err := row.Scan(&userID, &role)

	// If the user vas not found abort the transaction and return error
	if err == sql.ErrNoRows {
		t.Rollback()
		return LoginReturn{}, errors.New("username and password combination not found")
	}

	// generate a token
	token := uuid.New().String()
	_, err = t.Exec(`INSERT INTO TokenTable VALUES (UUID_TO_BIN(?),?,NOW());`, token, userID)
	if err != nil {
		t.Rollback()
		return LoginReturn{}, err
	}

	log.Printf("user %s logged in, generated token %s", userName, token)
	t.Commit()
	return LoginReturn{UserID: userID, Token: token, Role: role}, nil
}
