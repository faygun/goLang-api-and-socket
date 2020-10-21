package product

import (
	"database/sql"
	"log"

	"github.com/faygun/goLang-api-and-socket/database"
)

func getProduct(productID int) (*Product, error) {
	row := database.DbConn.QueryRow("SELECT manufacturer, pricePerUnit, productId, productName, quantityOnHand, sku, upc FROM products where productId = ? ", productID)
	product := Product{}
	err := row.Scan(&product.Manufacturer, &product.PricePerUnit, &product.ProductID, &product.ProductName, &product.QuantityOnHand, &product.Sku, &product.Upc)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &product, nil
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query("SELECT manufacturer, pricePerUnit, productId, productName, quantityOnHand, sku, upc FROM products")
	if err != nil {
		return nil, err
	}

	defer results.Close()
	list := make([]Product, 0)
	for results.Next() {
		product := Product{}
		results.Scan(&product.Manufacturer, &product.PricePerUnit, &product.ProductID, &product.ProductName, &product.QuantityOnHand, &product.Sku, &product.Upc)
		list = append(list, product)
	}

	return list, nil
}

func updateProduct(product Product) error {
	log.Print(1)
	_, err := database.DbConn.Exec(`UPDATE products set manufacturer=?, pricePerUnit=CAST(? AS DECIMAL(13,2)), productId=?, 
								productName=?, quantityOnHand=?, sku=?, upc=? WHERE productId =?`,
		product.Manufacturer, product.PricePerUnit, product.ProductID,
		product.ProductName, product.QuantityOnHand, product.Sku, product.Upc, product.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func insertProduct(product Product) (int, error) {
	result, err := database.DbConn.Exec(`INSERT INTO products (manufacturer, pricePerUnit, productName, quantityOnHand, sku, upc)
	VALUES(?, ?, ?, ?, ?, ?)`, product.Manufacturer, product.PricePerUnit,
		product.ProductName, product.QuantityOnHand, product.Sku, product.Upc)

	if err != nil {
		return 0, err
	}

	insertID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(insertID), nil
}

func deleteProduct(productID int) error {
	_, err := database.DbConn.Exec(`DELETE FROM products WHERE productId=?`, productID)

	if err != nil {
		return err
	}

	return nil
}
