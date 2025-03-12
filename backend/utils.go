package main

import (
	"fmt"
	"log"
	"net/http"
)

func sendAndLogError(w *http.ResponseWriter, errCode int, messages ...string) {
	(*w).WriteHeader(errCode)
	log.Print(messages)
	fmt.Fprint(*w, messages)
}
