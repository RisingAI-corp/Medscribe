package utils

import (
	"fmt"
	"net/http"
)

// GetAccessToken retrieves the access token from cookies.
func GetAccessToken(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("access token not found")
}

// GetRefreshToken retrieves the refresh token from cookies.
func GetRefreshToken(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("refresh token not found")
}
