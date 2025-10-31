package middleware

import(
	"net/http"
)

func (cfg *apiConfig) MiddlewareMain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		cfg.fileserverHits.Add(1)
	})
}
