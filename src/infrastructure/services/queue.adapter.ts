/**
 * QUEUE ADAPTER (BullMQ)
 *
 * Provides a base queue abstraction for async job processing.
 *
 * Responsibilities:
 * - Create named queues
 * - Add jobs to queues
 * - Configure retry strategies (exponential backoff)
 * - Set concurrency limits
 *
 * Each feature creates its own queue using this adapter.
 * Workers that process jobs live in shared/workers/.
 */
