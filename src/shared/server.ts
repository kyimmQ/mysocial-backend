/**
 * HTTP / SOCKET.IO SERVER
 *
 * Handles:
 * - Creating HTTP server from Express app
 * - Attaching Socket.IO to the HTTP server
 * - Configuring Socket.IO with Redis adapter
 * - Starting the server on configured port
 * - Graceful shutdown (SIGTERM, SIGINT)
 *
 * Separated from express.app.ts because:
 * - Express app = middleware + routes (testable without starting server)
 * - Server = actual listening + WebSocket lifecycle (integration concern)
 */
