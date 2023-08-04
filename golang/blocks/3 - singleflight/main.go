package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/sync/singleflight"
	"log"
	"net/http"
	"net/url"
)

var group singleflight.Group

func main() {
	conn, err := connectDB()
	if err != nil {
		fmt.Printf("connect db err : %v", err)
		return
	}

	app := App{
		db:      conn,
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("/", app.SlowResource)

	log.Fatal(http.ListenAndServe(":3000", app))
}

// App behaves kinda wrapper around DefaultServeMux to hold database connection.
type App struct {
	db *sql.DB
	http.Handler
}

func (a App) SlowResource(w http.ResponseWriter, r *http.Request) {

	// you need a key so that DO function realize what function is called. You can use the same group with different DO's.
	key := "key"
	_, err, shared := group.Do(key, func() (interface{}, error) {
		// let's say your database gives the response in 4 seconds.
		return a.db.Query("SELECT pg_sleep(4);")
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.Write([]byte(fmt.Sprintln("Is this shared response?", shared)))
}

// connecting to database
func connectDB() (*sql.DB, error) {
	username, password := "postgres", "postgres"

	q := make(url.Values)
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(username, password),
		Host:     "localhost",
		Path:     "singleflight",
		RawQuery: q.Encode(),
	}

	conn, err := sql.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping error : %w", err)
	}

	return conn, err
}
