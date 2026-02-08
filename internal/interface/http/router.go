// Package http contains the HTTP router configuration.
//
// Aggregates all feature handlers and mounts them under versioned paths.
//
// Example:
//
//	func NewRouter(container *app.Container) http.Handler {
//	    r := chi.NewRouter()
//
//	    // Global middleware
//	    r.Use(middleware.Recovery)
//	    r.Use(middleware.Logger)
//	    r.Use(middleware.CORS)
//
//	    // API v1
//	    r.Route("/api/v1", func(r chi.Router) {
//	        // Public
//	        r.Post("/signup", container.UserHandler.SignUp)
//	        r.Post("/signin", container.UserHandler.SignIn)
//
//	        // Protected
//	        r.Group(func(r chi.Router) {
//	            r.Use(middleware.Auth)
//	            r.Get("/currentuser", container.UserHandler.CurrentUser)
//	            r.Get("/users/{id}", container.UserHandler.GetProfile)
//	        })
//	    })
//
//	    return r
//	}
package http
