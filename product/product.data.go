package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/faygun/goLang-api-and-socket/database"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

// func init() {
// 	fmt.Println("loading...")
// 	prodMap, err := loadProductMap()
// 	productMap.m = prodMap
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%d products loaded...\n", len(productMap.m))
// }

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal(file, &productList)
	if err != nil {
		log.Fatal(err)
	}

	prodMap := make(map[int]Product)
	for _, product := range productList {
		prodMap[product.ProductID] = product
	}

	return prodMap, nil
}

func getProduct(productID int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[productID]; ok {
		return &product
	}

	return nil
}

func removeProduct(productID int) {
	productMap.RLock()
	defer productMap.RUnlock()
	delete(productMap.m, productID)
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query("SELECT * FROM products")
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

func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductID() int {
	productIds := getProductIds()
	return productIds[len(productIds)-1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	// if the product id is set, update, otherwise add
	addOrUpdateID := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		// if it exists, replace it, otherwise return error
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}
