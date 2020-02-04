package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/auth0-community/auth0"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	jose "gopkg.in/square/go-jose.v2"
)

// StatusHandler: return status
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//these should obvs be put in .env file
		secret := []byte("eTahGTp5PzjxLVczMlMLn3tuPvCJCclq")
		secretProvider := auth0.NewKeyProvider(secret)
		audience := []string{"authtest"}

		configuration := auth0.NewConfiguration(secretProvider, audience, "https://dev-45yrn41w.eu.auth0.com/", jose.HS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))

	r.Handle("/status", authMiddleware(StatusHandler)).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Our application will run on port 3000. Here we declare the port and pass in our router.
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}
