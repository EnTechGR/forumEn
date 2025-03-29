package handlers

import (
	"fmt"
	"forum/middleware"
	"net/http"
)

// HomeHandler handles requests to the home page ("/")
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	user := middleware.GetCurrentUser(r)

	if user == nil {
		fmt.Fprintln(w, "User is not logged in")
	} else {
		fmt.Fprintf(w, "Hello %s", user.Username)
	}
}
