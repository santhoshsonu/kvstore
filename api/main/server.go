package main

import (
	"fmt"
	"kvstore/api/handlers"
	"kvstore/store"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const PORT = 8080

const T_LOG_FILE = "./transaction.log"

func main() {
	err := store.InitializeStore(T_LOG_FILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	router := mux.NewRouter()

	url := "/api/v1/key/{key}"
	router.HandleFunc(url, handlers.GetHandler).Methods("GET")
	router.HandleFunc(url, handlers.PutHandler).Methods("PUT")
	router.HandleFunc(url, handlers.DeleteHandler).Methods("DELETE")

	log.Printf("Server Listening on PORT: %d\n", PORT)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), router))
}
