package restapi

import (
	"log"
	"net/http"
	"time"
	"os"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		// Logging into a file
		file, err := os.OpenFile("REST-API.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Println("Failed to open log file REST-API.log", err)
		}
		var MyFile *log.Logger
		MyFile = log.New(file,
    			"REST-API-logs: ",
    			log.Ldate|log.Ltime|log.Lshortfile)
	
		MyFile.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
