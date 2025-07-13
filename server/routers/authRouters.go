package routers

import (
	"net/http"
	"personal-diary/controllers"
	"personal-diary/middleware"

	"github.com/gorilla/mux"
)

func AuthRouters(routers *mux.Router) {
	// Implement the authentication routes here
	// routers.Use(middleware.CORSMiddleware)

	// Google OAuth routes
	routers.HandleFunc("/auth/google/login", controllers.GoogleLogin).Methods("GET")
	routers.HandleFunc("/auth/google/callback", controllers.GoogleCallback).Methods("GET")

	routers.HandleFunc("/register", controllers.Register).Methods("POST")
	routers.HandleFunc("/login", controllers.Login).Methods("POST")

	routers.HandleFunc("/forgot-password", controllers.RequestPasswordReset).Methods("POST") // or POST with body
	routers.HandleFunc("/validate-reset-token", controllers.ValidateResetToken).Methods("POST")
	routers.HandleFunc("/reset-password", controllers.ResetPassword).Methods("POST")



	routers.Handle("/protected", middleware.JwtVerify(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access granted to protected route"))
	}))).Methods("GET")
	routers.HandleFunc("/logout", controllers.Logout).Methods("POST")
	// whenever user go back to /login, /register or / automatically log out
}
