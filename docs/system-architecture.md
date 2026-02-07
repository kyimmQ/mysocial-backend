# MySocial Backend - System Architecture

## Architecture Overview

MySocial Backend follows a **layered architecture** with clear separation of concerns, asynchronous processing, and multi-layered caching. The system is designed for high scalability, performance, and real-time user interactions.

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
│              (Web Browser, Mobile App, REST Clients)             │
└───────────────────────┬─────────────────────────────────────────┘
                        │ HTTPS / WebSocket
                        ↓
┌─────────────────────────────────────────────────────────────────┐
│                    Load Balancer (AWS ALB)                       │
└───────────────────────┬─────────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        ↓               ↓               ↓
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   Node.js    │ │   Node.js    │ │   Node.js    │
│  Instance 1  │ │  Instance 2  │ │  Instance 5  │
│  (PM2 Cluster)│ │              │ │              │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
        ┌───────────────┼───────────────┐
        ↓               ↓               ↓
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│    Redis     │ │   MongoDB    │ │  Cloudinary  │
│   (Cache)    │ │ (Persistent) │ │   (Media)    │
└──────────────┘ └──────────────┘ └──────────────┘
        │
        ↓
┌──────────────┐
│  Bull Queues │
│  (Workers)   │
└──────────────┘
```

---

## System Components

### 1. Presentation Layer

**Components**:
- Express.js HTTP server
- Socket.IO WebSocket server
- Route handlers
- Controllers
- Middleware (auth, validation, error handling)

**Responsibilities**:
- Accept HTTP requests and WebSocket connections
- Validate request input using Joi schemas
- Authenticate users via JWT
- Route requests to appropriate controllers
- Return HTTP responses
- Emit real-time events to connected clients

**Technologies**: Express v4, Socket.IO v4, Joi v17

---

### 2. Business Logic Layer

**Components**:
- Feature modules (9 domains)
- Service classes
- Validation schemas
- Business logic

**Responsibilities**:
- Implement business rules
- Coordinate between cache, database, and queue layers
- Transform data between layers
- Generate computed values (JWT tokens, hashed passwords)

**Technologies**: TypeScript v4.9, Lodash v4

---

### 3. Data Access Layer

**Components**:
- Database services
- Mongoose models and schemas
- Query builders

**Responsibilities**:
- Abstract MongoDB operations
- Provide type-safe data access
- Handle database connections
- Execute queries and aggregations

**Technologies**: Mongoose v7, MongoDB v5+

---

### 4. Caching Layer

**Components**:
- Redis cache services
- Base cache abstraction
- Cache invalidation logic

**Responsibilities**:
- Store frequently accessed data in-memory
- Reduce database load
- Provide sub-millisecond read latency
- Support real-time features (sorted feeds, leaderboards)

**Technologies**: Redis v4+, ioredis client

---

### 5. Queue & Worker Layer

**Components**:
- Bull queues (12 queues)
- Queue workers (11 workers)
- Job processors

**Responsibilities**:
- Process heavy operations asynchronously
- Retry failed jobs with exponential backoff
- Ensure eventual consistency between cache and database
- Send emails asynchronously
- Process media uploads

**Technologies**: Bull v4, BullMQ v3, Bull Board v4

---

### 6. Real-time Communication Layer

**Components**:
- Socket.IO server
- Socket event handlers (6 types)
- Redis adapter for multi-instance support

**Responsibilities**:
- Maintain WebSocket connections
- Broadcast events to connected clients
- Handle room-based messaging
- Synchronize events across server instances

**Technologies**: Socket.IO v4, @socket.io/redis-adapter v8

---

### 7. External Services Layer

**Components**:
- Email service (SendGrid/Nodemailer)
- Media storage (Cloudinary)
- Third-party integrations

**Responsibilities**:
- Send transactional emails
- Upload and transform media files
- Integrate with external APIs

**Technologies**: SendGrid v7, Nodemailer v6, Cloudinary v1

---

## Request Flow Diagrams

### Synchronous Read Request (Cache-First)

```
Client
  │
  │ GET /api/users/:userId
  ↓
Express Route (/api/users/:userId)
  │
  │ authMiddleware validates JWT
  ↓
Controller (GetProfile)
  │
  │ 1. Check Redis cache
  ↓
Redis Cache
  │
  ├─── Cache HIT ────→ Return cached data ──→ Response (200ms)
  │
  └─── Cache MISS
         │
         │ 2. Query MongoDB
         ↓
       MongoDB
         │
         │ 3. Store in cache
         ↓
       Redis Cache
         │
         │ 4. Return data
         ↓
       Response (500ms)
```

---

### Asynchronous Write Request (Write-Through Cache)

```
Client
  │
  │ POST /api/posts
  ↓
Express Route (/api/posts)
  │
  │ authMiddleware validates JWT
  │ joiValidation validates body
  ↓
Controller (CreatePost)
  │
  │ 1. Save to Redis cache (immediate)
  ├──→ Redis Cache (writes in 10ms)
  │
  │ 2. Add job to Bull queue
  ├──→ Bull Queue (post queue)
  │      │
  │      │ 3. Worker processes job (async)
  │      ↓
  │    Queue Worker
  │      │
  │      │ 4. Write to MongoDB
  │      ↓
  │    MongoDB (persisted)
  │
  │ 5. Emit Socket.IO event
  ├──→ Socket.IO Server
  │      │
  │      │ 6. Broadcast to connected clients
  │      ↓
  │    Connected Clients (real-time update)
  │
  │ 7. Return success response immediately
  ↓
Response (150ms - doesn't wait for DB write)
```

---

### Real-time Messaging Flow

```
Sender Client
  │
  │ POST /api/chat/message
  ↓
Express Route (/api/chat/message)
  │
  │ authMiddleware validates JWT
  ↓
Controller (AddChatMessage)
  │
  │ 1. Save to Redis LIST (messages cache)
  ├──→ Redis Cache
  │
  │ 2. Add job to chat queue
  ├──→ Bull Queue (chat queue)
  │      │
  │      │ Worker writes to MongoDB
  │      ↓
  │    MongoDB (conversation + message)
  │
  │ 3. Emit Socket.IO event to receiver's room
  ├──→ Socket.IO Server
  │      │
  │      │ Check if receiver is online (in room)
  │      ├─── Online ──→ Emit 'message received' event
  │      │                  │
  │      │                  ↓
  │      │              Receiver Client (instant notification)
  │      │
  │      └─── Offline ──→ Queue notification job
  │                         │
  │                         ↓
  │                    Email notification sent later
  │
  │ 4. Update conversation list
  ├──→ Socket.IO Server
  │      │
  │      │ Emit 'chat list' to both sender and receiver
  │      ↓
  │    Both Clients (conversation updated)
  │
  │ 5. Return success response
  ↓
Response (100ms)
```

---

## Detailed Component Architecture

### Authentication Flow

```
1. Sign Up Flow
   Client → POST /signup
     │
     ├─→ SignUp Controller
     │     │
     │     ├─→ Validate input (Joi)
     │     ├─→ Hash password (bcrypt, 10 rounds)
     │     ├─→ Generate JWT token (userId, username, email)
     │     ├─→ Create session cookie (encrypted)
     │     ├─→ Save User to Redis HASH
     │     ├─→ Add job to Auth Queue
     │     │     │
     │     │     └─→ Auth Worker → Save to MongoDB (Auth + User collections)
     │     │
     │     └─→ Response { user, token }

2. Sign In Flow
   Client → POST /signin
     │
     ├─→ SignIn Controller
     │     │
     │     ├─→ Validate input (Joi)
     │     ├─→ Find user in MongoDB (Auth collection)
     │     ├─→ Compare password (bcrypt.compare)
     │     ├─→ Generate JWT token
     │     ├─→ Create session cookie
     │     ├─→ Load User from MongoDB → Save to Redis
     │     │
     │     └─→ Response { user, token }

3. Password Reset Flow
   Client → POST /forgot-password
     │
     ├─→ Password Controller
     │     │
     │     ├─→ Generate crypto token (20 bytes)
     │     ├─→ Save token to MongoDB (Auth.passwordResetToken)
     │     ├─→ Add job to Email Queue
     │     │     │
     │     │     └─→ Email Worker → Send reset email (SendGrid/Nodemailer)
     │     │
     │     └─→ Response { message: 'Email sent' }
     │
   Client → POST /reset-password/:token
     │
     ├─→ Password Controller
     │     │
     │     ├─→ Verify token (MongoDB lookup)
     │     ├─→ Hash new password (bcrypt)
     │     ├─→ Update password in MongoDB
     │     ├─→ Clear reset token
     │     ├─→ Add job to Email Queue (confirmation email)
     │     │
     │     └─→ Response { message: 'Password reset successful' }
```

---

### Post Creation & Distribution Flow

```
1. Create Post
   Client → POST /api/posts
     │
     ├─→ CreatePost Controller
     │     │
     │     ├─→ Validate input (Joi: text, image, video, gif, privacy)
     │     │
     │     ├─→ If image/video:
     │     │     └─→ CloudinaryUpload service
     │     │           └─→ Upload to Cloudinary (returns imgId, imgVersion)
     │     │
     │     ├─→ Create post object (postId, userId, text, media, reactions={}, commentsCount=0)
     │     │
     │     ├─→ Save to Redis (ZSET sorted by timestamp + HASH for post data)
     │     │     └─→ PostCache.savePostToCache(postId, postData)
     │     │
     │     ├─→ Add job to Post Queue
     │     │     └─→ Post Worker → MongoDB.create(postData)
     │     │
     │     ├─→ Emit Socket.IO event 'add post'
     │     │     └─→ All connected clients receive new post
     │     │
     │     └─→ Response { message: 'Post created', post }

2. Get Posts (Feed)
   Client → GET /api/posts?page=1
     │
     ├─→ GetPosts Controller
     │     │
     │     ├─→ Check Redis ZSET (sorted by timestamp, descending)
     │     │     │
     │     │     ├─── Cache HIT (posts in Redis)
     │     │     │     └─→ Return cached posts (10ms)
     │     │     │
     │     │     └─── Cache MISS
     │     │           │
     │     │           ├─→ Query MongoDB (PostModel.find().sort({ createdAt: -1 }))
     │     │           ├─→ Populate cache (PostCache.saveMultiplePostsToCache)
     │     │           └─→ Return posts (300ms)
     │     │
     │     └─→ Response { posts: [], totalPosts }
```

---

For detailed architecture information on specific subsystems, see:
- [Caching Architecture](./architecture/caching.md)
- [Queue Architecture](./architecture/queues.md)
- [Real-time Architecture](./architecture/realtime.md)
- [Database Design](./architecture/database.md)
