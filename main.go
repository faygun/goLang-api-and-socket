package main

import (
	"log"
	"net/http"

	"github.com/faygun/goLang-api-and-socket/receipt"

	"github.com/faygun/goLang-api-and-socket/database"

	"github.com/faygun/goLang-api-and-socket/product"
	_ "github.com/go-sql-driver/mysql"
)

const basePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(basePath)
	receipt.SetupRoutes(basePath)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
