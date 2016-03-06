package main

import (
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
)

// TODO : need to change to init ,this server it should run by core main.
func main() {
	fmt.Println("api server is started ... ", "\n", "port :8182 :443")
	router := NewRouter()
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServeTLS(":443", "/root/.lego/certificates/paas.jigko.net.crt", "/root/.lego/certificates/paas.jigko.net.key", handler))
	///root/.lego/certificates/paas.jigko.net.crt
}
