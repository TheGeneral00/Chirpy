package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/TheGeneral00/Chirpy/internal/helpers"
)

const maxBodyScan = 1 << 20 //1MiB

var keysToCheck = []string {
	"user", "username", "email", "password", "pass", "body", "message",
}

var headersToCheck = []string {
	"User", "Username", "X-User", "Authorization", "Proxy-Authorization",
}

func (cfg *APIConfig) InputSanatizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check headers 
		for _, hn := range headersToCheck {
			if hv := r.Header.Get(hn); hv != "" {
				if helpers.CheckPatterns(hv) {
					helpers.RespondWithError(w, http.StatusBadRequest, "Suspicious input detected", nil)
					return
				}
			}
		}

		query := r.URL.Query()
		for name, vals := range query {
			if interestingField(name){
				for _, value := range vals{
					if helpers.CheckPatterns(value){
						helpers.RespondWithError(w, http.StatusBadRequest, "Suspicious input detected", nil)
						return
					}
				}
			}
		}
		var bodyBytes []byte 
		if r.Body != nil {
			limited := io.LimitReader(r.Body, maxBodyScan+1)
			b, _ := io.ReadAll(limited)
			bodyBytes = b 
			r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyBytes), r.Body))
		}
	contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "multipart/form-data") {
			_ = r.ParseForm()
			for name, vals := range r.PostForm {
				if interestingField(name) {
					for _, val := range vals {
						if helpers.CheckPatterns(val) {
							helpers.RespondWithError(w, http.StatusBadRequest, "Suspicious input detected", nil)
							return
						}
					}
				}
			}
		}

		if strings.Contains(contentType, "application/json") && len(bodyBytes) > 0 {
			var decoded any
			if err := json.Unmarshal(bodyBytes, &decoded); err == nil {
				if m, ok := decoded.(map[string]any); ok {
					if checkMapForPatterns(m) {
						helpers.RespondWithError(w, http.StatusBadRequest, "Suspicious input detected", nil)
						return
					}
				}
			}
		} else if len(bodyBytes) > 0 {
			bodyStr := string(bodyBytes)
			for _, k := range keysToCheck {
				if strings.Contains(strings.ToLower(bodyStr), k) {
					if helpers.CheckPatterns(bodyStr) {
						helpers.RespondWithError(w, http.StatusBadRequest, "Suspicious input detected", nil)
						return
					}
				}
			}
		}

		if auth := r.Header.Get("Authorization"); auth != "" {
			if helpers.CheckPatterns(auth) {
				helpers.RespondWithError(w, http.StatusBadRequest, "Suspicious input detected", nil)
				return
			}
		}
	next.ServeHTTP(w, r)
	})
}

func interestingField(name string) bool {
	name = strings.ToLower(name)
	for _, key := range keysToCheck {
		if name == key || strings.Contains(name, key) {
			return true
		}
	}
	return false 
}

//Used to walk a json object
func checkMapForPatterns(m map[string]any) bool {
	for key, v := range m {
		if interestingField(key) {
			switch vv := v.(type) {
			case string:
				if helpers.CheckPatterns(vv){
					return true
				}
			case []any:
				for _, e := range vv {
					if string, ok := e.( string ); ok{
						if helpers.CheckPatterns(string){
							return true
						}
					}
				}
			case map[string]any:
				if checkMapForPatterns(vv) {
					return true
				}
			}
		} else {
			if nested, ok := v.(map[string]any); ok{
				if checkMapForPatterns(nested){
					return true
				}		
			}
		}
	}
	return false
}

