package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	ProductId      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func Init() {
	productsJson :=
		`[
		{
			"productId" : 1,
			"manufacturer" : "Sony",
			"sku" : "123",
			"upc" : "test",
			"pricePerUnit" : "10.00",
			"quantityOnHand" : 25,
			"productName" : "walkman"
		},
		{
			"productId" : 2,
			"manufacturer" : "Panasonic",
			"sku" : "456",
			"upc" : "test",
			"pricePerUnit" : "15.00",
			"quantityOnHand" : 20,
			"productName" : "player"
		},
		{
			"productId" : 3,
			"manufacturer" : "Sharp",
			"sku" : "789",
			"upc" : "test",
			"pricePerUnit" : "12.00",
			"quantityOnHand" : 22,
			"productName" : "dvd"
		}
	]`

	err := json.Unmarshal([]byte(productsJson), &productList)
	if err != nil {
		log.Fatal("Cannot load product list", err)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(productsJson)
		if err != nil {
			log.Printf("Error when returning product list", err)
		}
	case http.MethodPost:
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var newProduct Product
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductId != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductId = 4
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusAccepted)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	product, index := findProductById(productId)
	if product == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	switch r.Method {
	case http.MethodGet:
		prodJson, err := json.Marshal(product)
		if err != nil {
			log.Printf("Unable to serialize product", err)
		}
		w.Header().Set("Content-type", "application/json")
		_, err = w.Write(prodJson)
	case http.MethodPut:
		productBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		var updatedProduct Product
		err = json.Unmarshal(productBytes, &updatedProduct)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		if product.ProductId != productId {
			w.WriteHeader(http.StatusBadRequest)
		}

		productList[index] = updatedProduct
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func findProductById(productId int) (*Product, int) {
	return &productList[productId], productId
}

func middleHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before handler start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middleware finished; %s\n", time.Since(start))
	})
}

func main() {
	Init()

	productsListHandler := http.HandlerFunc(productsHandler)
	productsItemHandler := http.HandlerFunc(productHandler)
	http.Handle("/products", middleHandler(productsListHandler))
	http.Handle("/products/", middleHandler(productsItemHandler))

	http.ListenAndServe(":5000", nil)
}
