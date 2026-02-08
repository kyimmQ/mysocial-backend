/**
 * BACKGROUND JOB WORKERS
 *
 * Process jobs from Bull/BullMQ queues.
 *
 * Pattern:
 *   Queue adds job → Worker picks it up → Worker calls repository → Data saved to DB
 *
 * Workers:
 *   - auth.worker.ts         → Save auth credentials
 *   - user.worker.ts         → Update user profiles
 *   - post.worker.ts         → Post CRUD persistence
 *   - reaction.worker.ts     → Reaction persistence
 *   - comment.worker.ts      → Comment persistence
 *   - follower.worker.ts     → Follower operations
 *   - chat.worker.ts         → Message persistence
 *   - notification.worker.ts → Notification creation
 *   - email.worker.ts        → Send emails via mailer adapter
 *   - image.worker.ts        → Image metadata updates
 *
 * Workers should be THIN — delegate actual DB work to repositories.
 */
