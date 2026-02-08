/**
 * SOCKET.IO EVENT HANDLERS
 *
 * Real-time event handlers organized by domain:
 *   - post.socket.ts     → new post, update, delete
 *   - chat.socket.ts     → message sent, read, typing
 *   - follower.socket.ts → follow, unfollow
 *   - notification.socket.ts → new notification
 *   - user.socket.ts     → online/offline status
 *
 * Each handler registers listeners on the Socket.IO server.
 * Uses Redis adapter to sync events across PM2 instances.
 */
