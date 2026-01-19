package router

import (
	"net/http"
	"ticketapp/internal/handlers"
)

func Register(mux *http.ServeMux, auth *handlers.AuthHandler) {
	mux.HandleFunc("/auth/login", auth.Login)
	mux.HandleFunc("/auth/refresh", auth.Refresh)
}
