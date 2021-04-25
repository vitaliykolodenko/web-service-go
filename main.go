package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	ProductId int `json:"productId"`
	Manufacturer string `json:"manufacturer"`
	Sku string `json:"sku"`
	Upc string `json:"upc"`
	PricePerUnit string `json:"pricePerUnit"`
	QuantityOnHand int `json:"quantityOnHand"`
	ProductName string `json:"productName"`
}

var productList []Product

func Init(){
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

func productsHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(productsJson)
		if err != nil {
			log.Printf("Error when returning product list", err)
		}
	case http.MethodPost:
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var newProduct Product
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductId != 0{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductId = 4
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusAccepted)
	}
}

func main() {
	Init()

	http.HandleFunc("/products", productsHandler)

	http.ListenAndServe(":5000", nil)
}
