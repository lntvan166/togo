package main

import (
	"fmt"
	"lntvan166/togo/internal/config"
	"lntvan166/togo/internal/controller"
	"lntvan166/togo/internal/repository"
	"lntvan166/togo/internal/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	config.Load()

	db := repository.Connect()
	defer db.Close()

	controller.NewHandler(db)

	route := mux.NewRouter()
	routes.HandleRequest(route)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started!")

}
