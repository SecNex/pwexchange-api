package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get all headers
		start := time.Now()
		clientIP := r.RemoteAddr
		method := r.Method
		urlPath := r.URL.Path
		httpVersion := r.Proto
		userAgent := r.UserAgent()

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Milliseconds()
		log.Printf("%s - \"%s %s %s - %s\" %d bytes -> %s %d in %dms\n",
			clientIP,
			method,
			urlPath,
			httpVersion,
			userAgent,
			rw.size,
			http.StatusText(rw.statusCode),
			rw.statusCode,
			duration,
		)
	})
}
