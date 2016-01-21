package restapi

import (
	"log"
	"net/http"
)

func RunRestApi() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":6060", router))
}
