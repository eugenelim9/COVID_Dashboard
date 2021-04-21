package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//create the data source name, which identifies the
	//user, password, server address, and default database
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/demo", os.Getenv("MYSQL_ROOT_PASSWORD"))

	//create a database object, which manages a pool of
	//network connections to the database server
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	//ensure that the database gets closed when we are done
	defer db.Close()

	//for now, just ping the server to ensure we have
	//a live connection to it
	if err := db.Ping(); err != nil {
		fmt.Printf("error pinging database: %v\n", err)
	} else {
		fmt.Printf("successfully connected!\n")
	}
}
