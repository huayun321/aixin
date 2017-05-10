package main

import (
	"fmt"
	"github.com/urfave/negroni"
	"net/http"
	"github.com/gorilla/mux"
	"os"
)

var port = os.Getenv("PORT")

//test
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	fmt.Fprintln(w, "hello from heroku")
}

func main() {
	n := negroni.Classic()

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler).Methods("GET")

	n.UseHandler(router)
	n.Run(":" + port)
}
