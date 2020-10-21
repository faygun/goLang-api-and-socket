package database

import (
	"database/sql"
	"log"
	"time"
)

var DbConn *sql.DB

func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "root:password@tcp(192.168.99.101:32769)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}
	DbConn.SetConnMaxLifetime(60 * time.Second)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetMaxOpenConns(4)
}
