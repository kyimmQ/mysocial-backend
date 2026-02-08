// Package main is the composition root.
// It initializes all dependencies and starts the HTTP server.
//
// Flow:
//   1. Load config (infrastructure/config)
//   2. Connect to external services (database, cache, messaging)
//   3. Build the DI container (app/container.go)
//   4. Register routes (interface/http/router.go)
//   5. Start server
//
// No business logic here — only wiring.
package main
