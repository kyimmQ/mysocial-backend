/**
 * MONGODB CONNECTION LIFECYCLE
 *
 * Handles:
 * - Connection to MongoDB (mongoose.connect)
 * - Connection error handling
 * - Graceful disconnect on shutdown
 * - Connection event logging
 *
 * Pure infrastructure — no business models here.
 * Mongoose models live inside each feature's repositories/.
 */
