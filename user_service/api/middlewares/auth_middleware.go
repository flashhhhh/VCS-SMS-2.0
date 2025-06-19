package middlewares

import (
	"net/http"

	"github.com/flashhhhh/pkg/jwt"
	"github.com/flashhhhh/pkg/logging"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		logging.LogMessage("user_service", "Checking if header " + authHeader + " is admin", "DEBUG")
		if (authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.LogMessage("user_service", "Header " + authHeader + " doesn't contain token information", "INFO")
			return
		}

		token := authHeader[7:]

		logging.LogMessage("user_service", "Checking if token " + token + " is admin", "DEBUG")
		data, validateTokenErr := jwt.ValidateToken(token)
		if (validateTokenErr != nil) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.LogMessage("user_service", "Can't validate token " + token, "INFO")
			return
		}

		if (data["role"] != "admin") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			logging.LogMessage("user_service", "Not an admin token", "INFO")
			return 
		}

		logging.LogMessage("user_service", "This user is an admin! Forwarding to next handler.", "INFO")

		next.ServeHTTP(w, r)
	})
}

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		logging.LogMessage("user_service", "Checking if header " + authHeader + " is a user", "DEBUG")
		if (authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.LogMessage("user_service", "Header " + authHeader + " doesn't contain token information", "INFO")
			return 
		}

		token := authHeader[7:]
		
		logging.LogMessage("user_service", "Checking if token " + token + " is a user", "DEBUG")
		data, validateTokenErr := jwt.ValidateToken(token)
		if (validateTokenErr != nil) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.LogMessage("user_service", "Can't validate token " + token, "INFO")
			return 
		}

		if (data["role"] != "admin" && data["role"] != "user") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.LogMessage("user_service", "User is not an admin or a user", "INFO")
			return 
		}

		logging.LogMessage("user_service", "Forward this request to next handler, adding user's id","INFO")

		if data["role"] == "admin" {
			r.Header.Set("userRole", "admin")
		} else {
			r.Header.Set("userID", data["id"].(string))
		}
		next.ServeHTTP(w, r)
	})
}