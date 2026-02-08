/**
 * EXPRESS APPLICATION SETUP
 *
 * Configures the Express app with:
 * - Security middleware (helmet, hpp, cors)
 * - Body parsers (json, urlencoded)
 * - Cookie/session middleware
 * - Compression
 * - Global error handler
 *
 * Does NOT start the server — that's server.ts's job.
 * Does NOT define routes — those come from features.
 */
