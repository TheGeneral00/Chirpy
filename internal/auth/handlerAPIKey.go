package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
        ApiKey := strings.TrimLeft(headers.Get("Authorization"), "ApiKey ")
        if ApiKey == "" {
                return "", fmt.Errorf("Failed to retrieve ApiKey from header")
        }
        return ApiKey, nil
}
