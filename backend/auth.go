package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

// import "database/sql"

// This returns true if the token passes is valid and the token owner matches the resourceOwner
// Note currently that if the suer is an admin, access is always given
func AuthByToken(token string, resourceOwner int, db *sql.DB) (bool, error) {
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
		return false, err
	}
	// if user has a admin role
	if role != 0 {
		return true, nil
		// if the user own the requested resource
	} else if userID == resourceOwner {
		return true, nil
		// if not return false
	} else {
		return false, nil
	}

}

// This struct contains the return of CreateToken
type LoginReturn struct {
	UserID int
	Token  string
}

// If the username and password combination exist a string of the uuid userid is returned
func CreateToken(userName string, password string, db *sql.DB) (LoginReturn, error) {
	// TODO error handling
	// begin the transaction
	t, _ := db.Begin()
	// see if users exist
	row := t.QueryRow("SELECT userID FROM Users WHERE Username = ? AND `Password(Hash)` = ?;", userName, password)

	// read user id
	var userID int
	err := row.Scan(&userID)

	// If the user vas not found abort the transaction and return error
	if err == sql.ErrNoRows {
		t.Rollback() // TODO error handling
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
	return LoginReturn{UserID: userID, Token: token}, nil
	// test if user pswd exist
}
