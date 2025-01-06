package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ServeCouponVerse(r *mux.Router, port string) error {
	server := http.Server{
		Handler: r,
		Addr:    port,
	}
	err := server.ListenAndServe()
	return err
}
