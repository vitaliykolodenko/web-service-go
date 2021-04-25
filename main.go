package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type fooHandler struct {
	message string
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	_, err := w.Write([]byte(f.message))
	if err != nil {
		fmt.Printf("Some error happened %v \n", err)
	}
}

func barHandler(w http.ResponseWriter, r *http.Request){
	_, err := w.Write([]byte("Bar called"))
	if err != nil {
		fmt.Printf("Some error happened %v \n", err)
	}
}

func main() {
	http.Handle("/foo", &fooHandler{message: "hello"})
	http.HandleFunc("/bar", barHandler)

	jsonTest()

	http.ListenAndServe(":5000", nil)
}

func jsonTest(){
	type prod struct {
		Id int
		Name string
	}

	rez, err := json.Marshal(&prod{Id: 1, Name: "test"})
	if err != nil {
		fmt.Println("Got error when encoding json", err)
	}

	fmt.Printf("JSON: %v", string(rez))


	prodAfter := prod{}
	err = json.Unmarshal(rez, &prodAfter)

	if err != nil {
		fmt.Println("Got error when decoding json", err)
	}

	fmt.Printf("Product after decoding: %v", prodAfter)
}
