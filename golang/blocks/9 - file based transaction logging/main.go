package main

import (
	"github.com/fkaanoz/file-based-transaction-logger/internal/logger"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	fbLogger, err := logger.NewTransactionLogger("transaction.log")
	if err != nil {
		log.Fatal("logger creation error", err)
	}
	fbLogger.Run()

	a := Api{
		logger: fbLogger,
	}

	http.HandleFunc("/", a.TestHandler)

	log.Fatal(http.ListenAndServe(":4000", nil))
}

type Api struct {
	logger *logger.TransactionLogger
}

func (a *Api) TestHandler(w http.ResponseWriter, r *http.Request) {

	i := rand.Intn(10)
	if i < 5 {
		a.logger.WritePut(r.RemoteAddr, time.Now().String())
	} else {
		a.logger.WriteDelete(r.RemoteAddr)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test handler called"))
}
