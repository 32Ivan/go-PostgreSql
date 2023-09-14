package router

import (
	"go-PostgresSql/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/stacks/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stack", middleware.GetAllStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newstacks", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stack/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletestack/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")
	return router
}
