/**
 * COMPOSITION ROOT
 *
 * This is the entry point of the application.
 * It assembles all layers together:
 *
 * 1. Load and validate environment config
 * 2. Connect to external services (MongoDB, Redis)
 * 3. Initialize infrastructure adapters (mailer, queue)
 * 4. Wire up features (inject repositories into use-cases, use-cases into controllers)
 * 5. Register routes
 * 6. Start the HTTP/WebSocket server
 *
 * NO business logic here — only wiring.
 */
