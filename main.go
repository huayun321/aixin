package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"fmt"
	"os"
)
//test add mongo
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	fmt.Fprintf(w, "hello there")
}

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler).Methods("GET")
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" +port)
}