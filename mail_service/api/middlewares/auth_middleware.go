package middlewares

import (
	"net/http"

	"github.com/flashhhhh/pkg/jwt"
	"github.com/flashhhhh/pkg/logging"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		logging.LogMessage("mail_service", "Checking if header " + authHeader + " is admin", "DEBUG")
		if (authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.LogMessage("mail_service", "Header " + authHeader + " doesn't contain token information", "INFO")
			return
		}

		token := authHeader[7:]

		logging.LogMessage("mail_service", "Checking if token " + token + " is admin", "DEBUG")
		data, validateTokenErr := jwt.ValidateToken(token)
		if (validateTokenErr != nil) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.LogMessage("mail_service", "Can't validate token " + token, "INFO")
			return
		}

		if (data["role"] != "admin") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			logging.LogMessage("mail_service", "Not an admin token", "INFO")
			return 
		}

		logging.LogMessage("mail_service", "This user is an admin! Forwarding to next handler.", "INFO")

		next.ServeHTTP(w, r)
	})
}