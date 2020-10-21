package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/faygun/goLang-api-and-socket/database"
)

func getProduct(productID int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(ctx, "SELECT manufacturer, pricePerUnit, productId, productName, quantityOnHand, sku, upc FROM products where productId = ? ", productID)
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(ctx, "SELECT manufacturer, pricePerUnit, productId, productName, quantityOnHand, sku, upc FROM products")
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, `UPDATE products set manufacturer=?, pricePerUnit=CAST(? AS DECIMAL(13,2)), productId=?, 
								productName=?, quantityOnHand=?, sku=?, upc=? WHERE productId =?`,
		product.Manufacturer, product.PricePerUnit, product.ProductID,
		product.ProductName, product.QuantityOnHand, product.Sku, product.Upc, product.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func insertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products (manufacturer, pricePerUnit, productName, quantityOnHand, sku, upc)
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products WHERE productId=?`, productID)

	if err != nil {
		return err
	}

	return nil
}
