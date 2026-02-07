# MySocial Backend - Codebase Summary

## Overview

MySocial Backend is a comprehensive real-time social network application built with Node.js, TypeScript, Express, MongoDB, and Redis. The codebase follows a feature-based modular architecture with clear separation between business logic, data access, and presentation layers.

**Total Files**: 235 files
**Total Lines of Code**: ~117,835 tokens
**Primary Language**: TypeScript
**Test Coverage**: Jest unit tests for all controllers

## Directory Structure

```
mysocial-backend/
├── .circleci/               # CI/CD pipeline configuration
├── deployment/              # Terraform infrastructure as code
│   └── userdata/           # EC2 bootstrap scripts
├── endpoints/              # REST client test files (.http)
├── scripts/                # Deployment scripts (CodeDeploy)
├── src/                    # Main application source code
│   ├── features/           # Feature modules (9 domains)
│   ├── mocks/              # Test mock data
│   ├── shared/             # Shared services and utilities
│   │   ├── globals/        # Decorators and helpers
│   │   ├── services/       # Core services (DB, Redis, Email, Queues)
│   │   ├── sockets/        # Socket.IO event handlers
│   │   └── workers/        # Bull queue workers
│   ├── app.ts              # Application entry point
│   ├── config.ts           # Configuration management
│   ├── routes.ts           # Route aggregation
│   ├── seeds.ts            # Database seeding with Faker
│   ├── setupDatabase.ts    # MongoDB connection setup
│   └── setupServer.ts      # Express server configuration
└── [config files]          # ESLint, Prettier, Jest, TypeScript configs
```

## Feature Modules

The application is organized into 9 independent feature modules, each following a consistent structure:

### 1. Auth (`src/features/auth/`)
**Purpose**: User authentication and authorization

**Structure**:
- `controllers/` - SignUp, SignIn, SignOut, Password Reset, Current User
- `models/` - AuthModel (email, username, password, avatarColor, createdAt)
- `routes/` - authRoutes (public), currentRoutes (protected)
- `schemes/` - Joi validation schemas
- `interfaces/` - TypeScript type definitions

**Key Files**:
- `signup.ts` - Creates Auth + User records, generates JWT
- `signin.ts` - Validates credentials, creates session
- `password.ts` - Forgot/reset password with crypto tokens
- `current-user.ts` - Returns authenticated user profile

**Dependencies**: AuthService, UserCache, AuthQueue, EmailQueue

---

### 2. User (`src/features/user/`)
**Purpose**: User profile management and search

**Structure**:
- `controllers/` - Get Profile, Search Users, Change Password, Update Basic Info, Update Settings
- `models/` - UserModel (profilePicture, posts, followers, following, blocked, notifications, social links, etc.)
- `routes/` - userRoutes (protected), healthRoutes (public)
- `schemes/` - Joi validation for profile updates

**Key Files**:
- `get-profile.ts` - Cache-first profile retrieval, random user suggestions
- `search-user.ts` - MongoDB text search on username
- `update-basic-info.ts` - Quote, work, school, location
- `update-settings.ts` - Notification preferences

**Dependencies**: UserService, UserCache, UserQueue

---

### 3. Post (`src/features/post/`)
**Purpose**: Post creation, retrieval, updates, and deletion

**Structure**:
- `controllers/` - Create, Read, Update, Delete posts
- `models/` - PostModel (text, bgColor, privacy, imgVersion, imgId, videoVersion, videoId, feelings, gifUrl, reactions, commentsCount)
- `routes/` - postRoutes (protected)
- `schemes/` - Validation for post content, images, videos

**Key Files**:
- `create-post.ts` - Supports text, images, videos, GIFs (Cloudinary), background colors
- `get-posts.ts` - Pagination, sorting (newest first), cache-first
- `update-post.ts` - Partial updates with write-through cache
- `delete-post.ts` - Soft delete with cascade (removes reactions, comments)

**Dependencies**: PostService, PostCache, PostQueue, CloudinaryUpload

---

### 4. Reactions (`src/features/reactions/`)
**Purpose**: Post reactions (like, love, happy, wow, sad, angry)

**Structure**:
- `controllers/` - Add Reaction, Get Reactions, Remove Reaction
- `models/` - ReactionModel (postId, type, username, profilePicture)
- `routes/` - reactionRoutes (protected)
- `schemes/` - Reaction type validation

**Key Files**:
- `add-reactions.ts` - Upsert reaction, update post reaction count
- `get-reactions.ts` - Get all reactions for a post (grouped by type)
- `remove-reaction.ts` - Delete reaction, update counts

**Reaction Types**: `like`, `love`, `happy`, `wow`, `sad`, `angry`

**Dependencies**: ReactionService, ReactionCache, ReactionQueue

---

### 5. Comments (`src/features/comments/`)
**Purpose**: Comments on posts with notifications

**Structure**:
- `controllers/` - Add Comment, Get Comments
- `models/` - CommentModel (postId, comment, username, avatarColor, profilePicture, createdAt)
- `routes/` - commentRoutes (protected)
- `schemes/` - Comment validation

**Key Files**:
- `add-comment.ts` - Creates comment, increments post commentCount, sends notification, emits Socket.IO event
- `get-comments.ts` - Retrieves comments for a post with pagination

**Dependencies**: CommentService, CommentCache, CommentQueue, NotificationQueue

---

### 6. Followers (`src/features/followers/`)
**Purpose**: Follow/unfollow, block/unblock user relationships

**Structure**:
- `controllers/` - Follow User, Unfollow User, Get Followers, Block User
- `models/` - FollowerModel (followerId, followeeId, createdAt), BlockedUserModel
- `routes/` - followerRoutes (protected)

**Key Files**:
- `follower-user.ts` - Creates follower relationship, updates counts, sends notification
- `unfollow-user.ts` - Removes relationship, updates counts
- `get-followers.ts` - Get followers/following lists
- `block-user.ts` - Block/unblock users (prevents follows, messages)

**Dependencies**: FollowerService, FollowerCache, FollowerQueue, BlockedUserQueue

---

### 7. Chat (`src/features/chat/`)
**Purpose**: Real-time private messaging

**Structure**:
- `controllers/` - Add Message, Update Message, Delete Message, Get Messages, Add Message Reaction
- `models/` - MessageModel (conversationId, senderId, receiverId, body, gifUrl, isRead, deleteForMe, deleteForEveryone, reaction), ConversationModel (participants)
- `routes/` - chatRoutes (protected)
- `schemes/` - Message validation

**Key Files**:
- `add-chat-message.ts` - Sends message, updates conversation, emits Socket.IO event
- `update-chat-message.ts` - Mark as read
- `delete-chat-message.ts` - Soft delete (deleteForMe/deleteForEveryone)
- `add-message-reaction.ts` - React to messages

**Features**: Text, images, GIFs, message reactions, read receipts, soft delete

**Dependencies**: ChatService, MessageCache, ChatQueue, CloudinaryUpload

---

### 8. Images (`src/features/images/`)
**Purpose**: Image uploads and management via Cloudinary

**Structure**:
- `controllers/` - Add Image, Delete Image, Get Images
- `models/` - ImageModel (userId, imgVersion, imgId, createdAt)
- `routes/` - imageRoutes (protected)
- `schemes/` - Image validation

**Key Files**:
- `add-image.ts` - Upload to Cloudinary, save metadata, update user profile
- `delete-image.ts` - Delete from Cloudinary + database
- `get-images.ts` - Retrieve user images

**Dependencies**: ImageService, ImageQueue, CloudinaryUpload

---

### 9. Notifications (`src/features/notifications/`)
**Purpose**: In-app and email notifications

**Structure**:
- `controllers/` - Get Notifications, Update Notification, Delete Notification
- `models/` - NotificationModel (userTo, userFrom, message, notificationType, entityId, createdAt, read, comment, reaction, post, imgId, imgVersion, gifUrl)
- `routes/` - notificationRoutes (protected)

**Key Files**:
- `get-notifications.ts` - Retrieve user notifications with pagination
- `update-notification.ts` - Mark as read
- `delete-notification.ts` - Soft delete

**Notification Types**: `follows`, `comments`, `reactions`, `messages`

**Dependencies**: NotificationService, NotificationQueue, EmailQueue

---

## Shared Services

### Database Services (`src/shared/services/db/`)

All database services follow the same pattern:
- Abstract MongoDB queries
- Return promises
- Handle errors gracefully
- Use Mongoose models

**Services**:
- `auth.service.ts` - Auth CRUD operations
- `user.service.ts` - User profile operations, search, random users
- `post.service.ts` - Post CRUD, pagination, sorting
- `reaction.service.ts` - Reaction CRUD, grouping by type
- `comment.service.ts` - Comment CRUD
- `follower.service.ts` - Follower/following relationships
- `block-user.service.ts` - Block/unblock operations
- `chat.service.ts` - Message and conversation operations
- `image.service.ts` - Image metadata operations
- `notification.service.ts` - Notification CRUD

---

### Redis Cache Services (`src/shared/services/redis/`)

**Purpose**: High-performance caching layer to reduce database load

**Strategy**: Write-through cache (write to cache + queue to DB)

**Cache Services**:
- `user.cache.ts` - User profiles (HASH)
- `post.cache.ts` - Posts (HASH + ZSET for sorting)
- `reaction.cache.ts` - Reactions (HASH)
- `comment.cache.ts` - Comments (LIST)
- `follower.cache.ts` - Follower counts and lists (ZSET)
- `message.cache.ts` - Chat messages (LIST)

**Data Structures Used**:
- **HASH** - Store objects (users, posts, reactions)
- **ZSET** - Sorted sets (posts by date, followers)
- **LIST** - Ordered collections (messages, comments)
- **STRING** - Simple key-value (counts)

**Key Methods**:
- `saveToCache()` - Write to Redis
- `getFromCache()` - Read from Redis
- `deleteFromCache()` - Invalidate cache
- `updateCache()` - Partial updates

---

### Queue System (`src/shared/services/queues/`)

**Purpose**: Asynchronous job processing with Bull

**Queue Architecture**:
- Each queue handles 5 concurrent workers
- Jobs are retried 3 times on failure
- Exponential backoff retry strategy
- Bull Board UI for monitoring (`/queues`)

**Queues** (12 total):
- `auth.queue.ts` - User registration DB writes
- `user.queue.ts` - Profile updates
- `post.queue.ts` - Post CRUD operations
- `reaction.queue.ts` - Reaction DB writes
- `comment.queue.ts` - Comment DB writes
- `follower.queue.ts` - Follower relationship changes
- `blocked.queue.ts` - Block/unblock operations
- `chat.queue.ts` - Message persistence
- `image.queue.ts` - Image metadata updates
- `notification.queue.ts` - Notification creation
- `email.queue.ts` - Email sending (SendGrid/Nodemailer)

**Base Queue Pattern** (`base.queue.ts`):
```typescript
class BaseQueue {
  constructor(queueName: string);
  addJob(name: string, data: any): void;
  processJob(name: string, concurrency: number, callback: Function): void;
}
```

---

### Workers (`src/shared/workers/`)

**Purpose**: Process jobs from Bull queues

**Worker Pattern**:
- Listen to queue events
- Process jobs asynchronously
- Call database services
- Handle errors and logging

**Workers** (11 total):
- `auth.worker.ts` - Saves auth credentials to MongoDB
- `user.worker.ts` - Updates user profiles
- `post.worker.ts` - Post CRUD in MongoDB
- `reaction.worker.ts` - Reaction persistence
- `comment.worker.ts` - Comment persistence
- `follower.worker.ts` - Follower operations
- `blocked.worker.ts` - Block operations
- `chat.worker.ts` - Message persistence
- `image.worker.ts` - Image metadata updates
- `notification.worker.ts` - Notification creation
- `email.worker.ts` - Email sending

**Email Templates**:
- `forgot-password-template.ejs` - Password reset email
- `reset-password-template.ejs` - Password reset confirmation
- `notification.ejs` - In-app notification emails

---

### Socket.IO Handlers (`src/shared/sockets/`)

**Purpose**: Real-time bidirectional communication

**Socket Handlers**:
- `post.ts` - New post, update, delete events
- `follower.ts` - Follow/unfollow events
- `notification.ts` - Notification events
- `image.ts` - Image upload events
- `chat.ts` - Message events (send, read, typing)
- `user.ts` - User online/offline status

**Socket Events Emitted**:
- `add post` - New post created
- `update post` - Post updated
- `delete post` - Post deleted
- `add follower` - New follower
- `remove follower` - Unfollowed
- `insert notification` - New notification
- `update notification` - Notification read
- `message received` - New message
- `message read` - Message read
- `chat list` - Conversation list updated

**Redis Adapter**: Used for multi-instance Socket.IO synchronization

---

### Global Helpers (`src/shared/globals/`)

**Decorators** (`decorators/joi-validation.decorators.ts`):
- `@joiValidation(schema)` - Controller method validation decorator
- Validates request body against Joi schema
- Returns 400 Bad Request on validation failure

**Helpers** (`helpers/`):
- `auth-middleware.ts` - JWT verification, session validation
- `cloudinary-upload.ts` - Image/video upload to Cloudinary
- `error-handler.ts` - Global Express error handler
- `helpers.ts` - Utility functions (lowercase, capitalize, etc.)

**Error Classes**:
- `BadRequestError` - 400
- `NotFoundError` - 404
- `NotAuthorizedError` - 401
- `FileTooLargeError` - 413
- `ServerError` - 500

---

## Code Organization Patterns

### Feature Module Structure

Every feature follows this consistent pattern:

```
feature/
├── controllers/          # Request handlers
│   ├── test/            # Jest unit tests
│   └── *.ts             # Controller implementations
├── interfaces/          # TypeScript interfaces
├── models/              # Mongoose schemas
├── routes/              # Express route definitions
└── schemes/             # Joi validation schemas
```

### Controller Pattern

```typescript
class FeatureController {
  @joiValidation(schema)
  public async method(req: Request, res: Response): Promise<void> {
    // 1. Extract data from request
    // 2. Check cache (if applicable)
    // 3. Add job to queue (async write)
    // 4. Emit Socket.IO event (if real-time)
    // 5. Return response immediately
  }
}
```

### Service Layer Pattern

```typescript
class FeatureService {
  public async create(data: IData): Promise<void> {
    await Model.create(data);
  }

  public async getById(id: string): Promise<IData> {
    return await Model.findById(id);
  }

  public async update(id: string, data: Partial<IData>): Promise<void> {
    await Model.updateOne({ _id: id }, data);
  }

  public async delete(id: string): Promise<void> {
    await Model.deleteOne({ _id: id });
  }
}
```

### Cache-First Retrieval Pattern

```typescript
// 1. Check cache
let data = await cache.getFromCache(key);

// 2. If not in cache, fetch from DB
if (!data) {
  data = await service.getById(id);

  // 3. Populate cache
  await cache.saveToCache(key, data);
}

// 4. Return data
return data;
```

### Async Write Pattern

```typescript
// 1. Write to cache immediately
await cache.saveToCache(key, data);

// 2. Add job to queue for DB write
await queue.addJob('jobName', data);

// 3. Emit real-time event
socketIO.emit('eventName', data);

// 4. Return success response (don't wait for DB)
res.status(200).json({ message: 'Success' });
```

---

## File Naming Conventions

### Controllers
- Kebab-case with action prefix: `create-post.ts`, `get-profile.ts`, `add-comment.ts`
- Test files: `create-post.test.ts`

### Models
- Singular, kebab-case with `.schema` suffix: `user.schema.ts`, `post.schema.ts`

### Services
- Singular, kebab-case with `.service` suffix: `user.service.ts`, `post.service.ts`

### Queues
- Singular, kebab-case with `.queue` suffix: `post.queue.ts`, `email.queue.ts`

### Workers
- Singular, kebab-case with `.worker` suffix: `post.worker.ts`, `email.worker.ts`

### Routes
- Plural, camelCase with `Routes` suffix: `postRoutes.ts`, `userRoutes.ts`

### Interfaces
- Singular, kebab-case with `.interface` suffix: `post.interface.ts`

### Schemes (Joi)
- Plural or singular, kebab-case: `post.schemes.ts`, `signup.ts`

---

## Module Dependencies

### Core Dependencies
```
Express → Routes → Controllers → Services → MongoDB
                               → Cache → Redis
                               → Queue → Workers → Services
                               → Socket.IO
```

### Authentication Flow
```
Client → authRoutes → SignUpController → AuthService (create auth)
                                      → UserCache (save user)
                                      → AuthQueue (async DB write)
                                      → JWT (generate token)
                                      → Response (200 + token)
```

### Post Creation Flow
```
Client → postRoutes → CreatePostController → CloudinaryUpload (if media)
                                          → PostCache (save post)
                                          → PostQueue (async DB write)
                                          → Socket.IO (emit 'add post')
                                          → Response (201)
```

### Real-time Messaging Flow
```
Client → chatRoutes → AddChatMessageController → MessageCache (save)
                                               → ChatQueue (async DB)
                                               → Socket.IO (emit to receiver)
                                               → Response (200)
```

---

## Testing Strategy

### Unit Tests
- **Location**: `*/controllers/test/*.test.ts`
- **Framework**: Jest + ts-jest
- **Coverage**: All controllers have unit tests
- **Mocks**: Located in `src/mocks/` directory

### Mock Data
- `auth.mock.ts` - Auth and user mock data
- `post.mock.ts` - Post mock data
- `chat.mock.ts` - Message mock data
- `reactions.mock.ts` - Reaction mock data
- `followers.mock.ts` - Follower mock data
- `image.mock.ts` - Image mock data
- `notification.mock.ts` - Notification mock data
- `user.mock.ts` - User profile mock data

### Test Utilities
- Uses `@faker-js/faker` for generating realistic test data
- Mocks Redis, MongoDB, Socket.IO
- Tests request/response flow without actual DB/cache connections

---

## Configuration Files

### TypeScript (`tsconfig.json`)
- Target: ES2021
- Module: CommonJS
- Path aliases: `@auth/*`, `@user/*`, `@post/*`, `@shared/*`, etc.
- Strict mode enabled

### ESLint (`.eslintrc.json`)
- TypeScript ESLint parser
- Airbnb-style guide (modified)
- Prettier integration

### Prettier (`.prettierrc.json`)
- 2-space indentation
- Single quotes
- Semicolons required
- Trailing commas: ES5

### Jest (`jest.config.ts`)
- ts-jest preset
- Coverage threshold: 80%
- Test environment: node
- Transform: TypeScript files

---

## Build and Deployment

### Development
```bash
npm run dev          # Nodemon with tsconfig-paths
npm run test         # Jest with coverage
npm run lint:check   # ESLint
npm run prettier:fix # Format code
```

### Production Build
```bash
npm run build        # TypeScript compile + tsc-alias
npm start            # PM2 cluster (5 instances)
```

### Seeding
```bash
npm run seeds:dev    # Generate fake data with Faker
npm run seeds:prod   # Seed production DB
```

### PM2 Clustering
- 5 instances (configured in `package.json`)
- Auto-restart on file changes (watch mode)
- Log piping through Bunyan

---

## Summary Statistics

| Metric | Count |
|--------|-------|
| Feature Modules | 9 |
| Controllers | 35+ |
| Database Services | 10 |
| Cache Services | 6 |
| Queues | 12 |
| Workers | 11 |
| Socket Handlers | 6 |
| Mongoose Models | 11 |
| API Routes | 50+ |
| Unit Tests | 35+ |
| Total Files | 235 |

---

## Key Insights

1. **Separation of Concerns**: Clear boundaries between controllers, services, cache, and queues
2. **Async-First Design**: Non-blocking writes via Bull queues for better response times
3. **Cache-First Reads**: Redis reduces MongoDB load significantly
4. **Real-time Updates**: Socket.IO provides instant feedback to clients
5. **Horizontal Scalability**: PM2 clustering + Redis adapter enables multi-instance deployment
6. **Type Safety**: Full TypeScript coverage with strict mode
7. **Test Coverage**: Unit tests for all critical paths
8. **Infrastructure as Code**: Terraform for AWS deployment
9. **CI/CD Pipeline**: CircleCI with automated testing and deployment
