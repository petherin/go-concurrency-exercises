package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func slowQuery(ctx context.Context) error {
	//_, err := db.Exec("SELECT pg_sleep(5)")
	// Using ExecContext lets us pass the context so that
	// when the ctx.Done channel is closed, the db operation
	// terminates and returns error.
	_, err := db.ExecContext(ctx, "SELECT pg_sleep(5)")
	return err
}

func slowHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	err := slowQuery(req.Context()) // Sending req.Context() to slowQuery lets us use it to do cancellable calls
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
	fmt.Fprintln(w, "OK")
	fmt.Printf("slowHandler took: %v\n", time.Since(start))
}

func main() {
	// To run: go run main.go
	// In another terminal run time curl -i localhost:8000
	// Due to TimeoutHandler, we get the 'timeout!' message
	// and a 503 after 1 second.
	// The slowHandler will still take 5 seconds to complete though.
	// We need a way to propagate time awareness down to the handler function.
	// We can use req.Context() for this.

	var err error

	connstr := "host=localhost port=5432 user=alice password=pa$$word  dbname=wonderland sslmode=disable"

	db, err = sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}

	// This creates a timeout context which can be used in db.PingContext() to detect if
	// the db connection is still up. If ctx times out then we know the connection
	// is probably down and we get an error.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//if err = db.Ping(); err != nil {
	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	srv := http.Server{
		Addr:         "localhost:8000",
		WriteTimeout: 2 * time.Second,
		Handler: http.TimeoutHandler(http.HandlerFunc(slowHandler),
			1*time.Second,
			"timeout!"),
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}

// --> Installing postgres - macos
// brew install postgresql

// --> start
// pg_ctl -D /usr/local/var/postgres start

// --> create db and user
// psql postgres
// CREATE DATABASE wonderland;
// CREATE USER alice WITH ENCRYPTED PASSWORD 'pa$$word';
// GRANT ALL PRIVILEGES ON DATABASE wonderland TO alice;

// --> stop
// pg_ctl -D /usr/local/var/postgres stop

// --> postgresql download link
// https://www.postgresql.org/download/

// start postgresql - Windows
// pg_ctl -D "C:\Program Files\PostgreSQL\13\data" start

// stop postgresql - Windows
// pg_ctl -D "C:\Program Files\PostgreSQL\13\data" stop

// --> Linux
// sudo apt-get update
// sudo apt-get install postgresql-13

// sudo -u postgres psql -c "ALTER USER alice PASSWORD 'pa$$word';"
// sudo -u postgres psql -c "CREATE DATABASE wonderland;"

// sudo service postgresql start

// sudo service postgresql stop
