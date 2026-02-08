/**
 * REDIS CLIENT
 *
 * Creates and exports the Redis connection (ioredis).
 *
 * Used by:
 * - Feature cache services (via dependency injection)
 * - Bull/BullMQ queues (as connection config)
 * - Socket.IO Redis adapter (for multi-instance sync)
 */
