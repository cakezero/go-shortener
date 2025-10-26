package routes

import (
	"net/http"

	"github.com/cakezero/go-shortener/src/controllers"
	"github.com/cakezero/go-shortener/src/middlewares"
	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	// auth
	router.HandleFunc("/api/auth/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/auth/logout", controllers.Logout)
	router.HandleFunc("/api/auth/register", controllers.Register).Methods("POST")

	// server logic
	router.Handle("/api/delete-url", middlewares.AuthMiddleware(http.HandlerFunc(controllers.DeleteUrl))).Methods("DELETE")
	router.Handle("/api/delete-urls", middlewares.AuthMiddleware(http.HandlerFunc(controllers.DeleteAllUrls))).Methods("DELETE")
	router.Handle("/api/delete-selected-urls", middlewares.AuthMiddleware(http.HandlerFunc(controllers.DeleteSelectedUrls))).Methods("DELETE")
	router.Handle("/api/fetch-urls", middlewares.AuthMiddleware(http.HandlerFunc(controllers.FetchUrls)))
	router.HandleFunc("/api/shorten", controllers.Shorten).Methods("POST")
	router.HandleFunc("/api/visit-long-url", controllers.VisitLongUrl).Methods("POST")

	return router
}
