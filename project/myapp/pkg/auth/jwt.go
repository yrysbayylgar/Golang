package auth

import (
	"errors"
	"net/http"

	"myapp/internal/middleware"
)

// GetUserIDFromToken extracts the user ID from the token in the request context
func GetUserIDFromToken(r *http.Request) (string, error) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}