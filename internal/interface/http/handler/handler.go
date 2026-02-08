// Package handler contains HTTP handlers (controllers).
//
// Handlers are the bridge between HTTP and use-cases.
// They are THIN — no business logic allowed.
//
// Responsibilities:
//   1. Parse request (body, params, query, headers)
//   2. Validate input (delegate to pkg/validator)
//   3. Call use-case
//   4. Map domain errors to HTTP status codes
//   5. Write JSON response
//
// Example:
//
//	type UserHandler struct {
//	    createUser *user.CreateUserUseCase
//	    getUser    *user.GetUserUseCase
//	}
//
//	func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
//	    var input user.CreateUserInput
//	    json.NewDecoder(r.Body).Decode(&input)
//	    output, err := h.createUser.Execute(r.Context(), input)
//	    if err != nil {
//	        // map domain error → HTTP error
//	    }
//	    json.NewEncoder(w).Encode(output)
//	}
package handler
