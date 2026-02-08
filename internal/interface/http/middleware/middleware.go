// Package middleware provides HTTP middleware.
//
// Middleware functions wrap handlers to add cross-cutting behavior:
//   - AuthMiddleware: JWT verification, session validation
//   - LoggingMiddleware: request/response logging
//   - RecoveryMiddleware: panic recovery → 500 response
//   - CORSMiddleware: Cross-Origin Resource Sharing headers
//   - RateLimitMiddleware: per-IP request throttling
//
// Example:
//
//	func AuthMiddleware(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        token := r.Header.Get("Authorization")
//	        // verify JWT, attach user to context
//	        next.ServeHTTP(w, r)
//	    })
//	}
package middleware
