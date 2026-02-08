// Package domain contains domain-specific error types.
//
// These are business errors, not HTTP errors.
// The interface layer maps these to HTTP status codes.
//
// Examples:
//   - ErrUserNotFound      → handler maps to 404
//   - ErrDuplicateEmail    → handler maps to 409
//   - ErrInvalidCredentials → handler maps to 401
//   - ErrUnauthorized      → handler maps to 403
package domain
