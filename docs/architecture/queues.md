# Queue Architecture

## Overview

MySocial Backend uses **Bull** (Redis-based queue) for asynchronous job processing. This architecture enables non-blocking writes, reliable job execution, and horizontal scalability.

---

## Why Queues?

### Problems Solved

1. **Slow Database Writes**: MongoDB writes can take 100-500ms
2. **Blocking Operations**: Sending emails can take 1-3 seconds
3. **Media Processing**: Image/video uploads to Cloudinary take 500ms-2s
4. **User Experience**: Users shouldn't wait for slow operations

### Solution: Async Processing

```
Synchronous (Bad):
Client → Request → DB Write (500ms) → Email (2s) → Response
Total: 2.5 seconds

Asynchronous (Good):
Client → Request → Queue Job (10ms) → Response
Background: Worker → DB Write (500ms) → Email (2s)
Total: 10ms (from user perspective)
```

---

## Queue System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                        │
│  (Controllers add jobs to queues)                            │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ↓
┌─────────────────────────────────────────────────────────────┐
│                      Queue Layer                             │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐              │
│  │ Auth Queue │ │ Post Queue │ │ Email Queue│  ... (12)    │
│  └────────────┘ └────────────┘ └────────────┘              │
└───────────────────────┬─────────────────────────────────────┘
                        │ (stored in Redis)
                        ↓
┌─────────────────────────────────────────────────────────────┐
│                      Redis Storage                           │
│  (Job data, status, results)                                 │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ↓
┌─────────────────────────────────────────────────────────────┐
│                     Worker Layer                             │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐              │
│  │Auth Worker │ │Post Worker │ │Email Worker│  ... (11)    │
│  │(5 processes│ │(5 processes│ │(5 processes│              │
│  └────────────┘ └────────────┘ └────────────┘              │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ↓
┌─────────────────────────────────────────────────────────────┐
│                   External Services                          │
│  MongoDB, SendGrid, Cloudinary                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Queue List

### 1. AuthQueue (`auth.queue.ts`)

**Purpose**: Process user authentication operations

**Jobs**:
- `addAuthUserToDB` - Save new user credentials to MongoDB

**Worker**: `auth.worker.ts` (5 concurrent processes)

**Data Flow**:
```
SignUp Controller
  └─→ authQueue.addAuthUserJob('addAuthUserToDB', { value: authData })
        └─→ Auth Worker
              └─→ authService.createAuthUser(authData)
                    └─→ MongoDB (Auth collection)
```

---

### 2. UserQueue (`user.queue.ts`)

**Purpose**: Process user profile operations

**Jobs**:
- `addUserToDB` - Save new user profile to MongoDB
- `updateBasicInfoInDB` - Update user basic info
- `updateSocialLinksInDB` - Update social links
- `updateNotificationSettings` - Update notification preferences

**Worker**: `user.worker.ts` (5 concurrent processes)

**Data Flow**:
```
UpdateBasicInfo Controller
  └─→ userQueue.addUserJob('updateBasicInfoInDB', { userId, data })
        └─→ User Worker
              └─→ userService.updateBasicInfo(userId, data)
                    └─→ MongoDB (User collection)
```

---

### 3. PostQueue (`post.queue.ts`)

**Purpose**: Process post CRUD operations

**Jobs**:
- `addPostToDB` - Save new post to MongoDB
- `deletePostFromDB` - Delete post from MongoDB
- `updatePostInDB` - Update post in MongoDB

**Worker**: `post.worker.ts` (5 concurrent processes)

**Data Flow**:
```
CreatePost Controller
  └─→ postQueue.addPostJob('addPostToDB', { key: postId, value: postData })
        └─→ Post Worker
              └─→ postService.addPostToDB(postData)
                    └─→ MongoDB (Post collection)
```

---

### 4. ReactionQueue (`reaction.queue.ts`)

**Purpose**: Process post reactions

**Jobs**:
- `addReactionDataToDB` - Save reaction to MongoDB
- `removeReactionDataFromDB` - Remove reaction from MongoDB

**Worker**: `reaction.worker.ts` (5 concurrent processes)

**Data Flow**:
```
AddReaction Controller
  └─→ reactionQueue.addReactionJob('addReactionDataToDB', reactionData)
        └─→ Reaction Worker
              └─→ reactionService.addReactionDataToDB(reactionData)
                    └─→ MongoDB (Reaction collection)
```

---

### 5. CommentQueue (`comment.queue.ts`)

**Purpose**: Process post comments

**Jobs**:
- `addCommentToDB` - Save comment to MongoDB

**Worker**: `comment.worker.ts` (5 concurrent processes)

**Data Flow**:
```
AddComment Controller
  └─→ commentQueue.addCommentJob('addCommentToDB', commentData)
        └─→ Comment Worker
              └─→ commentService.addCommentToDB(commentData)
                    └─→ MongoDB (Comment collection)
```

---

### 6. FollowerQueue (`follower.queue.ts`)

**Purpose**: Process follower relationships

**Jobs**:
- `addFollowerToDB` - Create follower relationship
- `removeFollowerFromDB` - Remove follower relationship

**Worker**: `follower.worker.ts` (5 concurrent processes)

**Data Flow**:
```
FollowUser Controller
  └─→ followerQueue.addFollowerJob('addFollowerToDB', { followerId, followeeId })
        └─→ Follower Worker
              └─→ followerService.addFollowerToDB(followerId, followeeId)
                    └─→ MongoDB (Follower collection)
                    └─→ Update follower/following counts in User collection
```

---

### 7. BlockedQueue (`blocked.queue.ts`)

**Purpose**: Process user blocking operations

**Jobs**:
- `addBlockedUserToDB` - Block user
- `removeBlockedUserFromDB` - Unblock user

**Worker**: `blocked.worker.ts` (5 concurrent processes)

**Data Flow**:
```
BlockUser Controller
  └─→ blockedQueue.addBlockedUserJob('addBlockedUserToDB', { userId, blockedUserId })
        └─→ Blocked Worker
              └─→ blockUserService.blockUser(userId, blockedUserId)
                    └─→ MongoDB (User.blockedBy array)
```

---

### 8. ChatQueue (`chat.queue.ts`)

**Purpose**: Process chat messages and conversations

**Jobs**:
- `addChatMessageToDB` - Save message to MongoDB
- `markMessageAsDeleted` - Mark message as deleted
- `markMessagesAsReadInDB` - Mark messages as read
- `updateMessageReaction` - Add/remove message reaction

**Worker**: `chat.worker.ts` (5 concurrent processes)

**Data Flow**:
```
AddChatMessage Controller
  └─→ chatQueue.addChatJob('addChatMessageToDB', messageData)
        └─→ Chat Worker
              └─→ chatService.addMessageToDB(messageData)
                    └─→ MongoDB (Message collection)
                    └─→ Update conversation lastMessage
```

---

### 9. ImageQueue (`image.queue.ts`)

**Purpose**: Process image uploads and deletions

**Jobs**:
- `addImageToDB` - Save image metadata to MongoDB
- `updateBGImageInDB` - Update background image
- `addImageToProfile` - Add profile picture
- `removeImageFromDB` - Delete image metadata

**Worker**: `image.worker.ts` (5 concurrent processes)

**Data Flow**:
```
AddImage Controller
  ├─→ CloudinaryUpload (upload image, get imgId, imgVersion)
  └─→ imageQueue.addImageJob('addImageToDB', imageData)
        └─→ Image Worker
              └─→ imageService.addImage(imageData)
                    └─→ MongoDB (Image collection)
```

---

### 10. NotificationQueue (`notification.queue.ts`)

**Purpose**: Process in-app notifications

**Jobs**:
- `updateNotification` - Create or update notification
- `deleteNotification` - Delete notification

**Worker**: `notification.worker.ts` (5 concurrent processes)

**Data Flow**:
```
AddComment Controller (triggers notification)
  └─→ notificationQueue.addNotificationJob('updateNotification', notificationData)
        └─→ Notification Worker
              └─→ notificationService.insertNotification(notificationData)
                    └─→ MongoDB (Notification collection)
```

---

### 11. EmailQueue (`email.queue.ts`)

**Purpose**: Send emails asynchronously

**Jobs**:
- `forgotPasswordEmail` - Send password reset email
- `commentsEmail` - Send comment notification email
- `followersEmail` - Send follower notification email
- `reactionsEmail` - Send reaction notification email
- `directMessageEmail` - Send message notification email
- `changePassword` - Send password change confirmation

**Worker**: `email.worker.ts` (5 concurrent processes)

**Data Flow**:
```
ForgotPassword Controller
  └─→ emailQueue.addEmailJob('forgotPasswordEmail', { email, username, resetLink })
        └─→ Email Worker
              └─→ mailTransport.sendEmail(template, data)
                    ├─→ SendGrid (production)
                    └─→ Nodemailer (development)
```

---

## Base Queue Pattern

### BaseQueue Class

All queues extend the `BaseQueue` class:

```typescript
export abstract class BaseQueue {
  queue: Queue;
  log: Logger;

  constructor(queueName: string) {
    this.queue = new Queue(queueName, `${config.REDIS_HOST}`);
    this.log = config.createLogger(`${queueName}Queue`);

    // Event listeners
    this.queue.on('completed', (job: Job) => {
      job.remove();
    });

    this.queue.on('global:completed', (jobId: string) => {
      this.log.info(`Job ${jobId} completed`);
    });

    this.queue.on('global:stalled', (jobId: string) => {
      this.log.info(`Job ${jobId} is stalled`);
    });
  }

  protected addJob(name: string, data: IAuthJob | IEmailJob | IPostJobData): void {
    this.queue.add(name, data, {
      attempts: 3,
      backoff: { type: 'fixed', delay: 5000 }
    });
  }

  protected processJob(name: string, concurrency: number, callback: Queue.ProcessCallbackFunction<void>): void {
    this.queue.process(name, concurrency, callback);
  }
}
```

---

## Queue Configuration

### Job Options

```typescript
this.queue.add(jobName, jobData, {
  attempts: 3,                  // Retry up to 3 times
  backoff: {
    type: 'fixed',              // Fixed delay between retries
    delay: 5000                 // 5 seconds between retries
  },
  removeOnComplete: true,       // Auto-remove completed jobs
  removeOnFail: false           // Keep failed jobs for debugging
});
```

### Worker Concurrency

Each queue processes **5 jobs concurrently**:

```typescript
this.processJob('addAuthUserToDB', 5, authWorker.addAuthUserToDB);
```

**Total Concurrency**: 12 queues × 5 workers = **60 concurrent jobs**

---

## Worker Pattern

### Worker Implementation

```typescript
class AuthWorker {
  async addAuthUserToDB(job: Job, done: DoneCallback): Promise<void> {
    try {
      const { value } = job.data;

      // 1. Extract job data
      const authData: IAuthDocument = value;

      // 2. Perform database operation
      await authService.createAuthUser(authData);

      // 3. Update job progress
      job.progress(100);

      // 4. Mark job as done
      done(null, job.data);
    } catch (error) {
      // 5. Log error
      log.error(error);

      // 6. Mark job as failed
      done(error as Error);
    }
  }
}
```

### Error Handling

**Retry Logic**:
1. Job fails on first attempt
2. Wait 5 seconds (backoff delay)
3. Retry (attempt 2/3)
4. If fails again, wait 5 seconds
5. Retry (attempt 3/3)
6. If still fails, mark as failed permanently

**Failed Jobs**:
- Stored in Redis for inspection
- Can be manually retried via Bull Board UI
- Logged with full error details

---

## Job Data Structures

### Auth Job

```typescript
interface IAuthJob {
  value: IAuthDocument;
}

interface IAuthDocument {
  _id?: string | ObjectId;
  username: string;
  email: string;
  password?: string;
  avatarColor: string;
  createdAt?: Date;
}
```

---

### Post Job

```typescript
interface IPostJobData {
  key?: string;           // For updates/deletes
  value?: IPostDocument;  // For creates/updates
  keyOne?: string;        // Additional keys
  keyTwo?: string;
}

interface IPostDocument {
  _id?: string | ObjectId;
  userId: string;
  username: string;
  email: string;
  avatarColor: string;
  profilePicture: string;
  post: string;
  bgColor: string;
  imgVersion?: string;
  imgId?: string;
  videoVersion?: string;
  videoId?: string;
  feelings?: string;
  gifUrl?: string;
  privacy?: string;
  commentsCount: number;
  reactions: IReactions;
  createdAt?: Date;
}
```

---

### Email Job

```typescript
interface IEmailJob {
  receiverEmail: string;
  template: string;
  subject: string;
}

// Specific email types
interface IForgotPasswordParams {
  username: string;
  email: string;
  ipaddress: string;
  date: string;
  resetLink: string;
}

interface INotificationTemplate {
  username: string;
  message: string;
  header: string;
}
```

---

### Message Job

```typescript
interface IMessageData {
  _id?: string | ObjectId;
  conversationId: string;
  receiverId: string;
  receiverUsername: string;
  receiverAvatarColor: string;
  receiverProfilePicture: string;
  senderUsername: string;
  senderId: string;
  senderAvatarColor: string;
  senderProfilePicture: string;
  body: string;
  isRead?: boolean;
  gifUrl?: string;
  selectedImage?: string;
  reaction?: IMessageReaction[];
  createdAt?: Date;
  deleteForEveryone?: boolean;
  deleteForMe?: boolean;
}
```

---

## Bull Board UI

### Monitoring Dashboard

**Access**: `http://localhost:5000/queues`

**Features**:
- View all queues and their status
- See active, completed, failed, delayed jobs
- Inspect job data and results
- Manually retry failed jobs
- Clean completed/failed jobs
- Real-time updates

**Setup** (in `setupServer.ts`):
```typescript
import { ExpressAdapter } from '@bull-board/express';
import { createBullBoard } from '@bull-board/api';
import { BullAdapter } from '@bull-board/api/bullAdapter';

const serverAdapter = new ExpressAdapter();
serverAdapter.setBasePath('/queues');

createBullBoard({
  queues: [
    new BullAdapter(authQueue),
    new BullAdapter(userQueue),
    new BullAdapter(postQueue),
    // ... other queues
  ],
  serverAdapter
});

app.use('/queues', serverAdapter.getRouter());
```

---

## Performance Characteristics

### Throughput

**Single Worker**:
- Simple DB writes: 50-100 jobs/minute
- Email sending: 20-30 emails/minute
- Image processing: 10-20 uploads/minute

**5 Concurrent Workers**:
- Simple DB writes: 250-500 jobs/minute
- Email sending: 100-150 emails/minute
- Image processing: 50-100 uploads/minute

### Latency

**Job Processing Time**:
- Database writes: 100-500ms
- Email sending: 1-3 seconds
- Image uploads: 500ms-2 seconds

**Queue Overhead**: ~5-10ms (Redis communication)

---

## Job Lifecycle

```
1. Job Created
   ↓
2. Added to Queue (Redis)
   ↓
3. Worker picks up job
   ↓
4. Job processing starts
   ├─→ Success
   │     ├─→ Job marked as completed
   │     ├─→ Job removed from queue (if removeOnComplete: true)
   │     └─→ Event emitted: 'completed'
   │
   └─→ Failure
         ├─→ Job marked as failed (attempt 1/3)
         ├─→ Wait backoff delay (5 seconds)
         ├─→ Retry job (attempt 2/3)
         ├─→ If still failing, retry again (attempt 3/3)
         └─→ If all attempts fail:
               ├─→ Job marked as permanently failed
               ├─→ Event emitted: 'failed'
               └─→ Job kept in Redis for inspection
```

---

## Best Practices

### 1. Keep Job Data Small

```typescript
// Bad: Large job data
jobQueue.addJob('processUser', {
  user: entireUserObject,     // 5KB
  posts: allUserPosts,        // 50KB
  followers: allFollowers     // 100KB
});

// Good: Minimal job data
jobQueue.addJob('processUser', {
  userId: '507f1f77bcf86cd799439011'  // 24 bytes
});
// Worker fetches data from DB
```

---

### 2. Use Idempotent Operations

```typescript
// Jobs may be retried, so make operations idempotent
async addAuthUserToDB(job: Job): Promise<void> {
  const { value } = job.data;

  // Check if user already exists (idempotent)
  const existingUser = await AuthModel.findOne({ username: value.username });

  if (!existingUser) {
    await AuthModel.create(value);
  }
}
```

---

### 3. Handle Partial Failures

```typescript
async addPostToDB(job: Job): Promise<void> {
  const { value } = job.data;

  try {
    // Step 1: Create post
    await PostModel.create(value);

    // Step 2: Update user post count
    await UserModel.updateOne(
      { _id: value.userId },
      { $inc: { postsCount: 1 } }
    );
  } catch (error) {
    // If step 2 fails, rollback step 1
    await PostModel.deleteOne({ _id: value._id });
    throw error;
  }
}
```

---

### 4. Log Important Events

```typescript
async addAuthUserToDB(job: Job, done: DoneCallback): Promise<void> {
  log.info(`Processing job ${job.id}: addAuthUserToDB`);

  try {
    const { value } = job.data;
    await authService.createAuthUser(value);

    log.info(`Job ${job.id} completed successfully`);
    done(null, job.data);
  } catch (error) {
    log.error(`Job ${job.id} failed:`, error);
    done(error as Error);
  }
}
```

---

## Scaling Considerations

### Horizontal Scaling

**Current Setup**:
- 5 PM2 instances (each with own queue workers)
- Total: 5 instances × 12 queues × 5 workers = **300 concurrent jobs**

**Bottlenecks**:
1. **Redis**: Single Redis instance may become bottleneck
2. **MongoDB**: Write throughput may be limited
3. **External Services**: SendGrid/Cloudinary rate limits

**Solutions**:
1. **Redis Cluster**: Distribute queues across multiple Redis instances
2. **MongoDB Sharding**: Distribute data across multiple servers
3. **Rate Limiting**: Implement queue rate limits to respect external API limits

---

### Vertical Scaling

**Increase Worker Concurrency**:
```typescript
// From 5 to 10 concurrent workers per queue
this.processJob('addAuthUserToDB', 10, authWorker.addAuthUserToDB);
```

**Trade-offs**:
- ✅ Higher throughput
- ❌ Higher CPU/memory usage
- ❌ More database connections

---

## Monitoring & Debugging

### Key Metrics to Track

1. **Queue Length**: Number of waiting jobs
2. **Processing Rate**: Jobs processed per minute
3. **Failure Rate**: Percentage of failed jobs
4. **Average Processing Time**: Time per job
5. **Retry Rate**: Percentage of retried jobs

### Debugging Failed Jobs

**Via Bull Board**:
1. Navigate to `/queues`
2. Click on queue name
3. Click "Failed" tab
4. Inspect job data and error stack trace
5. Manually retry if needed

**Via Redis CLI**:
```bash
redis-cli LRANGE bull:queueName:failed 0 -1
```

---

## Future Improvements

1. **Priority Queues**: High-priority jobs (e.g., password resets) processed first
2. **Delayed Jobs**: Schedule jobs for future execution
3. **Job Chaining**: Link jobs together (job B starts when job A completes)
4. **Job Progress Tracking**: Real-time progress updates for long-running jobs
5. **Dead Letter Queue**: Separate queue for permanently failed jobs
6. **Rate Limiting**: Limit job processing rate to respect external API limits
7. **Job Metrics**: Collect and analyze job performance metrics
