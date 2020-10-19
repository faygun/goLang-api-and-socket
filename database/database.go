package database

import (
	"database/sql"
	"log"
)

var DbConn *sql.DB

func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "root:password@tcp(192.168.99.101:32769)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}
}
