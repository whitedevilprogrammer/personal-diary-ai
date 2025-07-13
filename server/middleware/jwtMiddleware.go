// package middleware

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/golang-jwt/jwt/v5"
// )

// var JwtKey = []byte("your_secret_key")
// var BlacklistedTokens = make(map[string]bool)

// // JwtVerify is a middleware function that checks for a valid JWT token in the Authorization header.
// // If the token is valid, it calls the next handler in the chain.
// func JwtVerify(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
// 		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 			return JwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// // RevokeToken is a function that revokes a JWT token by adding it to a blacklist.
// func RevokeToken(tokenString string) {
// 	// Implement token revocation logic here
// 	replaceToken := strings.Replace(tokenString, "Bearer ", "", 1)
// 	// For example, you could store the token in a database or cache
// 	cache := make(map[string]bool)
// 	cache[replaceToken] = true

// 	// This could involve adding the token to a blacklist in your database
// 	// or cache, depending on your application's architecture.
// }

// func IsBlacklisted(token string) bool {
// 	return BlacklistedTokens[token]
// }

package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var JwtKey = []byte("your_secret_key")
var BlacklistedTokens = make(map[string]bool)

// JwtMiddleware returns a Gorilla Mux middleware for JWT verification
func JwtMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Redirect to home page instead of returning an error
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

			// Check if token is blacklisted
			if IsBlacklisted(tokenStr) {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return JwtKey, nil
			})

			if err != nil || !token.Valid {
				// Redirect to home page instead of returning an error
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// For backward compatibility with existing code
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Redirect to login page instead of returning an error
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

		// Check if token is blacklisted
		if IsBlacklisted(tokenStr) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			// Redirect to home page instead of returning an error
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RevokeToken is a function that revokes a JWT token by adding it to the blacklist.
func RevokeToken(tokenString string) {
	// Remove "Bearer " prefix if present
	cleanToken := strings.Replace(tokenString, "Bearer ", "", 1)

	// Add token to the blacklist
	BlacklistedTokens[cleanToken] = true
}

// IsBlacklisted checks if a token is in the blacklist
func IsBlacklisted(token string) bool {
	return BlacklistedTokens[token]
}
