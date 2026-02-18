package middlewares

import (
	"log"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered: %v\n", r)
			}
		}()

		next.ServeHTTP(w, r)
	})

}
