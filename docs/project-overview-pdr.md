# MySocial Backend - Project Overview & Product Development Requirements

## Project Identity

**Project Name**: MySocial Backend (originally Chatty App Backend)

**Version**: 1.0.0

**Description**: A production-ready, full-featured real-time social media backend application that powers a Twitter/Facebook-like social network platform with comprehensive features including posts, reactions, comments, real-time messaging, notifications, and social connections.

**Repository**: Node.js/TypeScript backend server

**License**: ISC

---

## Executive Summary

MySocial Backend is an enterprise-grade social networking platform backend designed to support thousands of concurrent users with real-time interactions. The system leverages a modern microservices-inspired architecture with asynchronous job processing, multi-layered caching, and WebSocket-based real-time communication to deliver exceptional performance and scalability.

The platform handles the complete social media lifecycle including user authentication, content creation and curation, social graph management, private messaging, media uploads, and multi-channel notifications.

---

## Product Purpose

### Primary Objectives

1. **Provide a complete social networking backend** that handles all core features expected in modern social media platforms
2. **Demonstrate production-ready architecture** with scalability, performance, and reliability built-in from day one
3. **Enable real-time user interactions** through WebSocket connections and instant notification delivery
4. **Optimize for high performance** using multi-layered caching and asynchronous processing patterns
5. **Support cloud-native deployment** with infrastructure as code and CI/CD automation

### Problem Statement

Modern social media applications require complex backend infrastructure to handle:
- Real-time bidirectional communication between thousands of concurrent users
- High-throughput data processing for user-generated content
- Low-latency responses for feed generation and content discovery
- Reliable message delivery across multiple channels (in-app, email)
- Scalable media storage and delivery

MySocial Backend addresses these challenges through a carefully architected system that separates concerns, leverages caching strategically, and processes heavy workloads asynchronously.

---

## Target Audience

### Primary Users

1. **Full-stack Developers**
   - Need a production-ready backend for social media applications
   - Want to learn modern Node.js/TypeScript patterns
   - Require reference implementation for real-time features

2. **Backend Engineers**
   - Interested in scalable architecture patterns
   - Want to understand Redis caching strategies
   - Need examples of queue-based job processing

3. **DevOps Engineers**
   - Require infrastructure as code examples (Terraform)
   - Need CI/CD pipeline templates (CircleCI)
   - Want AWS deployment reference architecture

4. **Technical Architects**
   - Evaluating technology stack choices
   - Designing similar systems
   - Learning best practices for social platforms

### Secondary Users

- Students learning backend development
- Interview candidates preparing for system design questions
- Open-source contributors improving the codebase
- Startups building MVP social applications

---

## Key Features & Capabilities

### 1. Authentication & Authorization

**Requirements**:
- Users must be able to create accounts with email and password
- Users must be able to sign in and receive JWT tokens
- Users must be able to reset forgotten passwords via email
- Users must be able to change passwords while logged in
- System must maintain session state across requests

**Acceptance Criteria**:
- ✅ JWT-based authentication with secure token generation
- ✅ Bcrypt password hashing (10 salt rounds)
- ✅ Crypto-based password reset tokens (20-byte random)
- ✅ Cookie-based session management
- ✅ Protected routes requiring authentication
- ✅ Email verification for password resets

**Implementation**:
- Dual collection design (Auth + User) for separation of concerns
- Session stored in encrypted cookies
- Password reset tokens expire after configurable time
- JWT tokens include userId, username, email, avatarColor

---

### 2. User Profile Management

**Requirements**:
- Users must be able to view and edit their profiles
- Users must be able to upload profile pictures
- Users must be able to add social links (Instagram, Twitter, Facebook, YouTube)
- Users must be able to configure notification preferences
- Users must be able to search for other users by username
- System must suggest random users for discovery

**Acceptance Criteria**:
- ✅ Profile fields: basic info, quote, work, school, location
- ✅ Social links management
- ✅ Notification settings (messages, reactions, comments, follows)
- ✅ MongoDB text search on username field
- ✅ Random user suggestions using aggregation pipeline
- ✅ Cache-first profile retrieval for performance

**Implementation**:
- UserModel with embedded documents for social links and notifications
- Redis HASH for user profile caching
- Aggregation pipeline with `$sample` for random users
- Text index on username for fast search

---

### 3. Posts & Content

**Requirements**:
- Users must be able to create posts with text, images, videos, or GIFs
- Users must be able to set post privacy (public, private, followers-only)
- Users must be able to express feelings/activities with posts
- Users must be able to edit and delete their posts
- Users must be able to view paginated feeds
- System must support background colors for text posts

**Acceptance Criteria**:
- ✅ Support for text, image, video, and GIF posts
- ✅ Cloudinary integration for media storage
- ✅ Privacy settings (public, private, followers)
- ✅ Background color options for text posts
- ✅ Feelings/activity tags
- ✅ Pagination (default 10 posts per page)
- ✅ Sorted by creation date (newest first)
- ✅ Real-time updates via Socket.IO

**Implementation**:
- PostModel with flexible schema (optional imgId, videoId, gifUrl)
- Cloudinary upload helper for images/videos
- Redis ZSET for sorted post storage
- Socket events: `add post`, `update post`, `delete post`

---

### 4. Reactions

**Requirements**:
- Users must be able to react to posts with emotions
- System must support multiple reaction types
- Users must be able to change or remove reactions
- Post authors must see reaction counts
- Users must be able to view who reacted to posts

**Acceptance Criteria**:
- ✅ Six reaction types: like, love, happy, wow, sad, angry
- ✅ One reaction per user per post (upsert behavior)
- ✅ Real-time reaction count updates
- ✅ Reaction removal functionality
- ✅ Grouped reaction display by type
- ✅ Reaction notifications to post authors

**Implementation**:
- ReactionModel with type enum
- Redis HASH for reaction storage
- Atomic increment/decrement of post reaction counts
- Socket event: `update post` (with new reaction counts)

---

### 5. Comments

**Requirements**:
- Users must be able to comment on posts
- Users must be able to view all comments on a post
- Post authors must be notified of new comments
- System must track comment counts per post

**Acceptance Criteria**:
- ✅ Text-based comments with author info
- ✅ Pagination for comment lists
- ✅ Real-time comment notifications
- ✅ Comment count updates on posts
- ✅ Persistent storage in MongoDB
- ✅ Cache-first retrieval

**Implementation**:
- CommentModel with postId reference
- Redis LIST for comment caching
- Atomic increment of post commentsCount
- Socket event: `update post` (with new comment count)
- Notification queue job for comment notifications

---

### 6. Social Graph (Followers)

**Requirements**:
- Users must be able to follow and unfollow other users
- Users must be able to block and unblock users
- Users must see follower/following counts
- Blocked users must not appear in feeds or suggestions
- System must prevent circular blocks

**Acceptance Criteria**:
- ✅ Follow/unfollow functionality
- ✅ Block/unblock functionality
- ✅ Follower and following counts
- ✅ Blocked user list management
- ✅ Real-time follower notifications
- ✅ Privacy enforcement (blocked users)

**Implementation**:
- FollowerModel (followerId, followeeId)
- BlockedUserModel stored in User collection as array
- Redis ZSET for follower lists
- Atomic increment/decrement of counts
- Socket events: `add follower`, `remove follower`

---

### 7. Real-time Messaging

**Requirements**:
- Users must be able to send private messages to other users
- Messages must support text, images, and GIFs
- Users must be able to see read receipts
- Users must be able to react to messages
- Users must be able to delete messages (for me/for everyone)
- System must group messages into conversations
- Users must see online/offline status

**Acceptance Criteria**:
- ✅ Text, image, and GIF messages
- ✅ Read receipts (isRead flag)
- ✅ Message reactions
- ✅ Soft delete (deleteForMe, deleteForEveryone)
- ✅ Conversation grouping by participants
- ✅ Real-time message delivery via Socket.IO
- ✅ Message list pagination
- ✅ Typing indicators

**Implementation**:
- MessageModel with conversationId reference
- ConversationModel with participant array
- Redis LIST for message caching
- Socket events: `message received`, `message read`, `chat list`, `typing`
- Redis adapter for multi-instance Socket.IO

---

### 8. Image Management

**Requirements**:
- Users must be able to upload images to their profile
- System must support image gallery for each user
- Users must be able to delete uploaded images
- Images must be stored in cloud storage
- System must track image metadata

**Acceptance Criteria**:
- ✅ Cloudinary integration for storage
- ✅ Image metadata tracking (userId, imgVersion, imgId)
- ✅ Image deletion (Cloudinary + database)
- ✅ User image gallery retrieval
- ✅ Background image uploads for posts
- ✅ Profile picture updates

**Implementation**:
- ImageModel with Cloudinary identifiers
- CloudinaryUpload helper service
- Image queue for async processing
- Socket event: `add image`

---

### 9. Notifications

**Requirements**:
- Users must receive in-app notifications for relevant activities
- Users must receive email notifications based on preferences
- Users must be able to mark notifications as read
- Users must be able to delete notifications
- System must support multiple notification types

**Acceptance Criteria**:
- ✅ Notification types: follows, comments, reactions, messages
- ✅ In-app notification delivery
- ✅ Email notification delivery (configurable)
- ✅ Mark as read functionality
- ✅ Soft delete
- ✅ Pagination for notification list
- ✅ Real-time notification updates via Socket.IO

**Implementation**:
- NotificationModel with polymorphic references (comment, reaction, post)
- NotificationQueue for async creation
- EmailQueue for async email sending
- Socket event: `insert notification`, `update notification`
- SendGrid (production) and Nodemailer (development) for emails

---

## Technology Stack

### Backend Framework

**Node.js v16+**
- **Rationale**: Industry-standard for high-performance, event-driven JavaScript runtime
- **Benefits**: Non-blocking I/O, large ecosystem (npm), excellent for real-time applications

**TypeScript v4.9+**
- **Rationale**: Type safety reduces runtime errors, improves developer experience
- **Benefits**: Better IDE support, self-documenting code, easier refactoring

**Express.js v4**
- **Rationale**: Lightweight, flexible, and widely adopted web framework
- **Benefits**: Minimal overhead, excellent middleware ecosystem, easy to test

---

### Databases

**MongoDB v7**
- **Rationale**: Document-oriented NoSQL database ideal for social media data
- **Benefits**: Flexible schema, horizontal scaling, rich query language
- **Use Cases**: User profiles, posts, comments, messages, notifications

**Redis v4+**
- **Rationale**: In-memory data structure store for caching and real-time features
- **Benefits**: Sub-millisecond latency, rich data types (HASH, ZSET, LIST), pub/sub
- **Use Cases**: Session storage, caching, real-time feeds, leaderboards

---

### Real-time Communication

**Socket.IO v4**
- **Rationale**: WebSocket library with fallback mechanisms and reconnection logic
- **Benefits**: Real-time bidirectional communication, room-based broadcasting, Redis adapter for horizontal scaling
- **Use Cases**: Live notifications, chat, typing indicators, online presence

**@socket.io/redis-adapter v8**
- **Rationale**: Enables Socket.IO to work across multiple server instances
- **Benefits**: Horizontal scalability, session stickiness not required

---

### Async Job Processing

**Bull v4**
- **Rationale**: Redis-based queue for delayed and scheduled jobs
- **Benefits**: Reliable job processing, retries, prioritization, rate limiting
- **Use Cases**: Database writes, email sending, image processing

**BullMQ v3**
- **Rationale**: Modern rewrite of Bull with TypeScript support
- **Benefits**: Better performance, improved API, TypeScript-first

**@bull-board/express v4**
- **Rationale**: Web UI for monitoring Bull queues
- **Benefits**: Real-time queue monitoring, job inspection, manual job retry

---

### Media Storage

**Cloudinary**
- **Rationale**: Cloud-based image and video management service
- **Benefits**: Automatic optimization, responsive images, CDN delivery, transformations
- **Use Cases**: Profile pictures, post images/videos, chat media

---

### Email Services

**SendGrid (Production)**
- **Rationale**: Reliable transactional email service with high deliverability
- **Benefits**: 99% deliverability, analytics, template management, API

**Nodemailer (Development)**
- **Rationale**: SMTP email sender for local testing
- **Benefits**: Works with Ethereal.email for testing, no cost, easy setup

---

### Authentication & Security

**jsonwebtoken v9**
- **Rationale**: Industry-standard JWT implementation for stateless auth
- **Benefits**: Compact, URL-safe, verifiable without database lookup

**bcryptjs v2**
- **Rationale**: Password hashing library with adaptive cost factor
- **Benefits**: Resistant to brute-force, salted hashing, future-proof

**cookie-session v2**
- **Rationale**: Encrypted cookie-based session storage
- **Benefits**: Stateless sessions, no server-side storage, secure

**helmet v6**
- **Rationale**: Security middleware setting HTTP headers
- **Benefits**: XSS protection, clickjacking prevention, content type sniffing protection

**hpp (HTTP Parameter Pollution)**
- **Rationale**: Protects against parameter pollution attacks
- **Benefits**: Prevents array-based attacks on query parameters

---

### Validation

**Joi v17**
- **Rationale**: Schema validation library for request bodies
- **Benefits**: Declarative, type coercion, detailed error messages

**Custom Validation Decorators**
- **Rationale**: TypeScript decorators for controller method validation
- **Benefits**: DRY principle, separation of concerns, reusable

---

### Testing

**Jest v29**
- **Rationale**: Full-featured testing framework with built-in mocking
- **Benefits**: Zero config, snapshot testing, coverage reporting

**ts-jest v29**
- **Rationale**: TypeScript preprocessor for Jest
- **Benefits**: Test TypeScript without compilation step

**@faker-js/faker v7**
- **Rationale**: Generate realistic fake data for testing
- **Benefits**: Deterministic seeds, localization, rich API

---

### Build & Development

**ts-node v10**
- **Rationale**: TypeScript execution engine for Node.js
- **Benefits**: No compilation step during development

**tsc-alias v1**
- **Rationale**: Replaces path aliases in compiled output
- **Benefits**: Clean import paths (`@auth/*` instead of `../../../`)

**nodemon**
- **Rationale**: Auto-restart server on file changes
- **Benefits**: Fast development cycle

**ESLint v8**
- **Rationale**: Linting tool for TypeScript/JavaScript
- **Benefits**: Code consistency, early error detection

**Prettier v2**
- **Rationale**: Opinionated code formatter
- **Benefits**: Consistent formatting, no bikeshedding

---

### Deployment & Infrastructure

**PM2**
- **Rationale**: Production process manager for Node.js
- **Benefits**: Clustering (5 instances), auto-restart, log management, zero-downtime reloads

**Terraform**
- **Rationale**: Infrastructure as code for AWS resources
- **Benefits**: Version-controlled infrastructure, repeatable deployments, multi-environment support

**AWS Services**:
- **EC2** - Application servers
- **Auto Scaling Group** - Horizontal scaling
- **Application Load Balancer** - Traffic distribution
- **ElastiCache (Redis)** - Managed Redis cluster
- **Route53** - DNS management
- **Certificate Manager** - SSL/TLS certificates
- **S3** - Static file storage (env files, backups)
- **CodeDeploy** - Automated deployments
- **VPC** - Network isolation
- **CloudWatch** - Monitoring and logs

**CircleCI**
- **Rationale**: CI/CD pipeline automation
- **Benefits**: Automated testing, linting, deployment, Slack notifications

---

## System Requirements

### Development Environment

**Required**:
- Node.js 16.x or higher
- MongoDB 5.0+ (local or Atlas)
- Redis 7.0+ (local or cloud)
- TypeScript 4.9+
- npm 8+

**Optional**:
- Docker (for containerized Redis/MongoDB)
- MongoDB Compass (GUI for MongoDB)
- Redis Commander (GUI for Redis)
- Postman or REST Client (for API testing)

---

### Production Environment

**Infrastructure**:
- AWS account with IAM permissions
- Domain name (for Route53 hosted zone)
- Terraform v1.0+
- AWS CLI configured

**Third-party Services**:
- Cloudinary account (API key, secret, cloud name)
- SendGrid account (API key, verified sender)
- CircleCI account (for CI/CD)
- CodeCov account (for coverage reporting)

**Environment Variables** (see `.env.development.example`):
```
DATABASE_URL          # MongoDB connection string
JWT_TOKEN             # JWT secret
NODE_ENV              # development | production
SECRET_KEY_ONE        # Session encryption key 1
SECRET_KEY_TWO        # Session encryption key 2
CLIENT_URL            # Frontend application URL
API_URL               # Backend API URL
REDIS_HOST            # Redis host address
CLOUD_NAME            # Cloudinary cloud name
CLOUD_API_KEY         # Cloudinary API key
CLOUD_API_SECRET      # Cloudinary API secret
SENDER_EMAIL          # Nodemailer email (dev)
SENDER_EMAIL_PASSWORD # Nodemailer password (dev)
SENDGRID_API_KEY      # SendGrid API key (prod)
SENDGRID_SENDER       # SendGrid sender email
EC2_URL               # AWS metadata URL
```

---

## Success Criteria

### Functional Requirements

- ✅ **Authentication**: Users can sign up, sign in, reset passwords
- ✅ **User Profiles**: Users can view/edit profiles, search users
- ✅ **Posts**: Users can create/read/update/delete posts with media
- ✅ **Reactions**: Users can react to posts (6 types)
- ✅ **Comments**: Users can comment on posts
- ✅ **Followers**: Users can follow/unfollow, block/unblock
- ✅ **Messaging**: Users can send private messages with media
- ✅ **Images**: Users can upload/delete images via Cloudinary
- ✅ **Notifications**: Users receive in-app and email notifications

---

### Non-Functional Requirements

**Performance**:
- ✅ API response time < 200ms (cached reads)
- ✅ API response time < 500ms (database writes via queue)
- ✅ WebSocket message delivery < 100ms

**Scalability**:
- ✅ Horizontal scaling via PM2 clustering (5 instances)
- ✅ Redis adapter for Socket.IO multi-instance support
- ✅ Queue-based async processing (12 queues, 5 workers each)

**Reliability**:
- ✅ Job retry mechanism (3 retries with exponential backoff)
- ✅ Error handling with custom error classes
- ✅ Graceful shutdown handling

**Security**:
- ✅ JWT-based authentication
- ✅ Bcrypt password hashing (10 rounds)
- ✅ CORS protection
- ✅ Helmet security headers
- ✅ HPP protection
- ✅ Input validation (Joi schemas)

**Maintainability**:
- ✅ TypeScript strict mode
- ✅ ESLint + Prettier code quality
- ✅ Unit test coverage (controllers)
- ✅ Consistent code structure (feature modules)
- ✅ Path aliases for clean imports

**DevOps**:
- ✅ Infrastructure as code (Terraform)
- ✅ CI/CD pipeline (CircleCI)
- ✅ Automated testing in pipeline
- ✅ Code coverage reporting (CodeCov)
- ✅ Slack notifications for build status

---

## Technical Constraints

### Hard Constraints

1. **Database**: MongoDB and Redis (no SQL databases)
2. **Language**: Node.js with TypeScript
3. **Cloud Provider**: AWS (Terraform-managed)
4. **Message Queue**: Bull/BullMQ (Redis-based)
5. **Real-time**: Socket.IO (WebSocket protocol)

### Soft Constraints

1. **Response Time**: Target < 500ms for all API endpoints
2. **Concurrent Users**: Designed for 1000+ concurrent connections
3. **Job Processing**: Max 5 concurrent workers per queue
4. **Session Storage**: Cookie-based (no server-side sessions)
5. **Media Storage**: Cloudinary (no local file storage)

---

## Future Roadmap Considerations

### Short-term (Next 3-6 months)

1. **Stories/Reels Feature**
   - Ephemeral content (24-hour expiry)
   - Video-first posts
   - Viewer tracking

2. **Advanced Search**
   - Full-text search on post content
   - Hashtag support
   - User search filters (location, interests)

3. **Group Features**
   - Create/join groups
   - Group posts and discussions
   - Group admin roles

4. **Enhanced Analytics**
   - Post reach and engagement metrics
   - User activity heatmaps
   - Trending topics

### Mid-term (6-12 months)

1. **Recommendation Engine**
   - ML-based post recommendations
   - User suggestions based on interests
   - Trending content discovery

2. **Content Moderation**
   - AI-powered image/text moderation
   - User reporting system
   - Admin moderation dashboard

3. **Video Processing**
   - Video transcoding
   - Thumbnail generation
   - Live streaming support

4. **Mobile Push Notifications**
   - FCM integration
   - Push notification preferences
   - Silent notifications for background sync

### Long-term (12+ months)

1. **Microservices Migration**
   - Separate services for Auth, Posts, Chat, Notifications
   - API Gateway pattern
   - Service mesh (Istio/Linkerd)

2. **GraphQL API**
   - Alternative to REST API
   - Real-time subscriptions
   - Flexible queries

3. **Multi-region Deployment**
   - CDN integration
   - Geo-distributed databases
   - Edge computing for latency reduction

4. **Advanced Monetization**
   - Ads platform
   - Premium subscriptions
   - Creator monetization tools

---

## Known Limitations

1. **No API Rate Limiting**: Currently no built-in rate limiting (should add)
2. **No Pagination Cursors**: Using offset-based pagination (should migrate to cursor-based)
3. **No Database Transactions**: Some operations lack transaction support
4. **Limited Error Logging**: Should integrate structured logging (Datadog, New Relic)
5. **No API Versioning**: Single API version (should implement v1, v2 strategy)
6. **Hard-coded Configurations**: Some configs are hard-coded (should externalize)

---

## Compliance & Standards

### Code Quality

- ESLint rules enforced
- Prettier formatting enforced
- TypeScript strict mode enabled
- Unit test coverage tracked

### Security

- OWASP Top 10 considerations
- Regular dependency updates (npm audit)
- Secrets management (environment variables)
- HTTPS enforced in production

### Documentation

- Code comments for complex logic
- API endpoint documentation (`.http` files)
- README with setup instructions
- Terraform infrastructure documentation

---

## Conclusion

MySocial Backend represents a production-ready, scalable social media platform backend with comprehensive features and best-practice architecture patterns. The system is designed for high performance, reliability, and maintainability, making it suitable for both educational purposes and as a foundation for real-world social applications.
