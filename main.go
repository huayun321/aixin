package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"fmt"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	fmt.Fprintf(w, "Welcome to the home page!")
}

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler).Methods("GET")
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" +port)
}