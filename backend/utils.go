package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var simpleLogger *log.Logger

func sendAndLogError(w *http.ResponseWriter, errCode int, messages ...string) {
	if simpleLogger == nil {
		simpleLogger = log.New(os.Stderr, "", 0)
	}
	now := time.Now()
	(*w).WriteHeader(errCode)
	pc, filename, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("%s [error] in %s [%s:%d]: %s", now.String(), runtime.FuncForPC(pc).Name(), filename, line, strings.Join(messages, ""))
	simpleLogger.Print(s)
	fmt.Fprint(*w, s)
}
