package router

import (
	"net/http"

	"ticketapp/internal/handlers"
	"ticketapp/internal/middlewares"
	"ticketapp/internal/services"
)

// NewRouter builds and returns the application's HTTP router
func NewRouter(
	authHandler *handlers.AuthHandler,
	adminHandler *handlers.AdminHandler,
	jwtService *services.JWTService,
) http.Handler {

	mux := http.NewServeMux()

	// -------------------------
	// AUTH ROUTES (PUBLIC)
	// -------------------------

	mux.Handle(
		"/auth/login",
		middlewares.SecurityHeaders(
			middlewares.RateLimit(
				http.HandlerFunc(authHandler.Login),
			),
		),
	)

	mux.Handle(
		"/auth/verify-otp",
		middlewares.SecurityHeaders(
			middlewares.RateLimit(
				http.HandlerFunc(authHandler.VerifyOTP),
			),
		),
	)

	mux.Handle(
		"/auth/refresh",
		middlewares.SecurityHeaders(
			http.HandlerFunc(authHandler.Refresh),
		),
	)

	// mux.Handle(
	// 	"/auth/logout",
	// 	middlewares.SecurityHeaders(
	// 		http.HandlerFunc(authHandler.Logout),
	// 	),
	// )

	mux.Handle(
		"/auth/forgot-password",
		middlewares.SecurityHeaders(
			middlewares.RateLimit(
				http.HandlerFunc(authHandler.ForgotPassword),
			),
		),
	)

	mux.Handle(
		"/auth/reset-password",
		middlewares.SecurityHeaders(
			middlewares.RateLimit(
				http.HandlerFunc(authHandler.ResetPassword),
			),
		),
	)

	// -------------------------
	// 2FA SETUP (AUTH REQUIRED)
	// -------------------------

	// mux.Handle(
	// 	"/auth/setup-2fa",
	// 	middlewares.SecurityHeaders(
	// 		middlewares.AuthMiddleware(jwtService)(
	// 			http.HandlerFunc(authHandler.Setup2FA),
	// 		),
	// 	),
	// )

	// mux.Handle(
	// 	"/auth/confirm-2fa",
	// 	middlewares.SecurityHeaders(
	// 		middlewares.RateLimit(
	// 			http.HandlerFunc(authHandler.Confirm2FA),
	// 		),
	// 	),
	// )

	// -------------------------
	// ADMIN ROUTES (RBAC)
	// -------------------------

	mux.Handle(
		"/admin/users",
		middlewares.SecurityHeaders(
			middlewares.AuthMiddleware(jwtService)(
				middlewares.RequireRole("admin")(
					http.HandlerFunc(adminHandler.CreateUser),
				),
			),
		),
	)

	// -------------------------
	// HEALTH CHECK
	// -------------------------

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return mux
}
