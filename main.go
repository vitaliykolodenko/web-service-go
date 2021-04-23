package main

import (
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
	http.ListenAndServe(":5000", nil)
}
