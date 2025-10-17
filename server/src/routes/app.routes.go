package routes

import (
	"net/http"

	"github.com/cakezero/go-shortener/src/controllers"
	"github.com/cakezero/go-shortener/src/middlewares"
	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.Handle("/", middlewares.AuthMiddleware(http.HandlerFunc(controllers.Home)))
	router.Handle("/delete-url", middlewares.AuthMiddleware(http.HandlerFunc(controllers.DeleteUrl))).Methods("DELETE")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/logout", controllers.Logout)
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/shorten", controllers.Shorten).Methods("POST")
	router.HandleFunc("/visit-long-url", controllers.VisitLongUrl).Methods("POST")

	return router
}
