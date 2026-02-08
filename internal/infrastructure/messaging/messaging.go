// Package messaging handles message queue connections and job dispatching.
//
// Options:
//   - Asynq (Redis-based, similar to BullMQ) — recommended for this project
//   - RabbitMQ (AMQP)
//   - Kafka (event streaming)
//
// Provides:
//   - Queue client initialization
//   - Job enqueuing
//   - Worker registration and processing
//   - Retry strategies
//
// Example with Asynq:
//
//	type TaskDistributor interface {
//	    DistributeTaskSendEmail(ctx context.Context, payload *SendEmailPayload) error
//	    DistributeTaskSaveUser(ctx context.Context, payload *SaveUserPayload) error
//	}
//
//	type TaskProcessor interface {
//	    ProcessTaskSendEmail(ctx context.Context, task *asynq.Task) error
//	    ProcessTaskSaveUser(ctx context.Context, task *asynq.Task) error
//	    Start() error
//	}
package messaging
