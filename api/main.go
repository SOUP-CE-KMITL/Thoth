package main

import (
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
)

// TODO : need to change to init ,this server it should run by core main.
func main() {
	fmt.Println("api server is started ... ", "\n", "port :8182")
	router := NewRouter()
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8182", handler))
}
