# Caching Architecture

## Overview

MySocial Backend uses **Redis** as a high-performance caching layer to reduce database load and improve response times. The caching strategy follows a **cache-first read, write-through write** pattern.

---

## Caching Strategy

### Cache-First Read Pattern

```
1. Check cache first
   ↓
2. If found (cache hit):
   └─→ Return data immediately (< 10ms)

3. If not found (cache miss):
   ├─→ Query database (200-500ms)
   ├─→ Store result in cache
   └─→ Return data
```

### Write-Through Write Pattern

```
1. Write to cache immediately
   ├─→ Redis write (< 10ms)
   └─→ Client receives success response

2. Queue job for database write (async)
   └─→ Worker processes job eventually
       └─→ MongoDB write (100-500ms)
```

**Benefits**:
- **Fast reads**: Sub-10ms response times for cached data
- **Fast writes**: Don't block on database writes
- **Reduced DB load**: 70-90% reduction in database queries
- **Eventual consistency**: Database catches up asynchronously

---

## Redis Data Structures

### 1. HASH (Hash Maps)

**Use Case**: Store object data (users, posts, reactions, comments)

**Example - User Cache**:
```
Key: users:${userId}
Value: {
  _id: "507f1f77bcf86cd799439011",
  username: "john_doe",
  email: "john@example.com",
  avatarColor: "#9c27b0",
  postsCount: "42",
  followersCount: "150",
  followingCount: "200",
  createdAt: "2024-01-15T10:30:00.000Z"
}
```

**Operations**:
```typescript
// Save user to cache
await redis.HSET('users:123', {
  username: 'john_doe',
  email: 'john@example.com',
  avatarColor: '#9c27b0'
});

// Get user from cache
const user = await redis.HGETALL('users:123');

// Update specific field
await redis.HSET('users:123', 'postsCount', '43');

// Delete user from cache
await redis.DEL('users:123');
```

---

### 2. ZSET (Sorted Sets)

**Use Case**: Store ordered collections (posts sorted by timestamp, followers sorted by userId)

**Example - Post Feed**:
```
Key: post
Members: [
  { value: 'post:507f1f77bcf86cd799439011', score: 1674645000000 },
  { value: 'post:507f1f77bcf86cd799439012', score: 1674644000000 },
  { value: 'post:507f1f77bcf86cd799439013', score: 1674643000000 }
]
```

**Operations**:
```typescript
// Add post to sorted set (score = timestamp)
await redis.ZADD('post', {
  score: Date.now(),
  value: `post:${postId}`
});

// Get latest posts (descending order)
const postIds = await redis.ZRANGE('post', 0, 9, { REV: true });

// Get posts count
const count = await redis.ZCARD('post');

// Remove post from sorted set
await redis.ZREM('post', `post:${postId}`);
```

**Example - Follower List**:
```
Key: followers:${userId}
Members: [
  { value: 'follower:${followerId1}', score: ${followerId1} },
  { value: 'follower:${followerId2}', score: ${followerId2} }
]
```

---

### 3. LIST (Ordered Lists)

**Use Case**: Store sequential data (chat messages, comments)

**Example - Message List**:
```
Key: messages:${conversationId}
Values: [
  '{"messageId":"123","senderId":"user1","body":"Hello","timestamp":1674645000000}',
  '{"messageId":"124","senderId":"user2","body":"Hi","timestamp":1674646000000}',
  '{"messageId":"125","senderId":"user1","body":"How are you?","timestamp":1674647000000}'
]
```

**Operations**:
```typescript
// Add message to list (append to end)
await redis.RPUSH(`messages:${conversationId}`, JSON.stringify(messageData));

// Get messages (range)
const messages = await redis.LRANGE(`messages:${conversationId}`, 0, 49); // Get first 50

// Get list length
const count = await redis.LLEN(`messages:${conversationId}`);

// Trim list (keep only recent messages)
await redis.LTRIM(`messages:${conversationId}`, 0, 99); // Keep only 100 messages
```

---

### 4. STRING (Simple Key-Value)

**Use Case**: Store counts, flags, simple values

**Example - Follower Count**:
```
Key: followers:count:${userId}
Value: "150"
```

**Operations**:
```typescript
// Set count
await redis.SET(`followers:count:${userId}`, '150');

// Get count
const count = await redis.GET(`followers:count:${userId}`);

// Increment count
await redis.INCR(`followers:count:${userId}`);

// Decrement count
await redis.DECR(`followers:count:${userId}`);
```

---

## Cache Services

### UserCache

**Data Structure**: HASH + ZSET

**Cached Data**:
- User profile (HASH: `users:${userId}`)
- User list (ZSET: `user`)

**Methods**:
```typescript
class UserCache extends BaseCache {
  // Save user to cache
  async saveUserToCache(key: string, userUId: string, user: IUserDocument): Promise<void>

  // Get user from cache
  async getUserFromCache(userId: string): Promise<IUserDocument | null>

  // Update user in cache
  async updateSingleUserItemInCache(userId: string, prop: string, value: string): Promise<void>

  // Get total users count
  async getTotalUsersInCache(): Promise<number>

  // Get random users
  async getRandomUsersFromCache(userId: string, count: number): Promise<IUserDocument[]>
}
```

**Cache Keys**:
- `users:${userId}` - User HASH
- `user` - ZSET of all user IDs

---

### PostCache

**Data Structure**: HASH + ZSET

**Cached Data**:
- Post data (HASH: `posts:${postId}`)
- Post list sorted by timestamp (ZSET: `post`)

**Methods**:
```typescript
class PostCache extends BaseCache {
  // Save post to cache
  async savePostToCache(data: ISavePostToCache): Promise<void>

  // Get posts from cache
  async getPostsFromCache(key: string, start: number, end: number): Promise<IPostDocument[]>

  // Get total posts count
  async getTotalPostsInCache(): Promise<number>

  // Get posts with images
  async getPostsWithImagesFromCache(key: string, start: number, end: number): Promise<IPostDocument[]>

  // Get user posts
  async getUserPostsFromCache(key: string, userId: string): Promise<IPostDocument[]>

  // Update post in cache
  async updatePostInCache(key: string, updatedPost: IPostDocument): Promise<void>

  // Delete post from cache
  async deletePostFromCache(key: string, currentUserId: string): Promise<void>
}
```

**Cache Keys**:
- `posts:${postId}` - Post HASH
- `post` - ZSET of all post IDs sorted by timestamp

---

### MessageCache

**Data Structure**: LIST + HASH

**Cached Data**:
- Message list (LIST: `messages:${conversationId}`)
- Conversation data (HASH: `conversations:${conversationId}`)

**Methods**:
```typescript
class MessageCache extends BaseCache {
  // Add message to cache
  async addChatListToCache(senderId: string, receiverId: string, conversationId: string): Promise<void>

  // Add message to conversation
  async addChatMessageToCache(conversationId: string, message: IMessageData): Promise<void>

  // Get messages from cache
  async getChatMessagesFromCache(senderId: string, receiverId: string): Promise<IMessageData[]>

  // Mark messages as read
  async markMessageAsReadInCache(senderId: string, receiverId: string): Promise<void>

  // Update message in cache
  async updateChatMessagesInCache(senderId: string, receiverId: string, messageId: string, reaction: string): Promise<void>

  // Delete message from cache
  async markMessageAsDeletedInCache(senderId: string, receiverId: string, messageId: string, type: string): Promise<void>
}
```

**Cache Keys**:
- `messages:${conversationId}` - LIST of messages
- `conversations:${conversationId}` - HASH of conversation metadata

---

### ReactionCache

**Data Structure**: HASH

**Cached Data**:
- Reaction data (HASH: `reactions:${reactionId}`)

**Methods**:
```typescript
class ReactionCache extends BaseCache {
  // Save reaction to cache
  async savePostReactionToCache(key: string, reaction: IReactionDocument, postReactions: IReactions, type: string, previousReaction: string): Promise<void>

  // Remove reaction from cache
  async removePostReactionFromCache(key: string, reaction: IReactionDocument, postReactions: IReactions): Promise<void>

  // Get reactions from cache
  async getReactionsFromCache(postId: string): Promise<[IReactionDocument[], number]>

  // Get single reaction by username
  async getSingleReactionByUsernameFromCache(postId: string, username: string): Promise<[IReactionDocument, number] | []>
}
```

**Cache Keys**:
- `reactions:${postId}:${username}` - Reaction HASH
- Reaction counts stored in post HASH

---

### CommentCache

**Data Structure**: LIST + HASH

**Cached Data**:
- Comment list (LIST: `comments:${postId}`)
- Comment data (HASH: `comments:${commentId}`)

**Methods**:
```typescript
class CommentCache extends BaseCache {
  // Save comment to cache
  async savePostCommentToCache(postId: string, comment: string): Promise<void>

  // Get comments from cache
  async getCommentsFromCache(postId: string): Promise<ICommentDocument[]>

  // Get comments count
  async getCommentsNamesFromCache(postId: string): Promise<ICommentNameList[]>

  // Get single comment
  async getSingleCommentFromCache(postId: string, commentId: string): Promise<ICommentDocument[]>
}
```

**Cache Keys**:
- `comments:${postId}` - LIST of comment IDs
- `comments:${commentId}` - HASH of comment data

---

### FollowerCache

**Data Structure**: ZSET

**Cached Data**:
- Followers list (ZSET: `followers:${userId}`)
- Following list (ZSET: `following:${userId}`)

**Methods**:
```typescript
class FollowerCache extends BaseCache {
  // Save follower to cache
  async saveFollowerToCache(key: string, value: string): Promise<void>

  // Remove follower from cache
  async removeFollowerFromCache(key: string, value: string): Promise<void>

  // Update follower count in user cache
  async updateFollowersCountInCache(userId: string, prop: string, value: number): Promise<void>

  // Get followers from cache
  async getFollowersFromCache(key: string): Promise<IFollowerData[]>
}
```

**Cache Keys**:
- `followers:${userId}` - ZSET of follower IDs
- `following:${userId}` - ZSET of following IDs
- Counts stored in user HASH

---

## Cache Invalidation Strategies

### 1. Write-Through Invalidation

**When**: User updates data
**Action**: Update cache immediately, queue DB write

```typescript
// Update user profile
async updateUserProfile(userId: string, updates: Partial<IUserDocument>): Promise<void> {
  // 1. Update cache
  await userCache.updateSingleUserItemInCache(userId, 'quote', updates.quote);

  // 2. Queue DB update
  await userQueue.addUserJob('updateUserInDB', { userId, updates });

  // 3. Return success (don't wait for DB)
  return;
}
```

---

### 2. Delete-on-Write Invalidation

**When**: Data deleted
**Action**: Delete from cache immediately, queue DB delete

```typescript
// Delete post
async deletePost(postId: string, userId: string): Promise<void> {
  // 1. Delete from cache
  await postCache.deletePostFromCache(postId, userId);

  // 2. Queue DB delete
  await postQueue.addPostJob('deletePostFromDB', { postId });

  // 3. Emit Socket.IO event
  socketIO.emit('delete post', { postId });

  // 4. Return success
  return;
}
```

---

### 3. TTL-Based Expiration

**When**: Data changes infrequently
**Action**: Set expiration time on cache keys

```typescript
// Cache with TTL (Time To Live)
await redis.SETEX(`session:${sessionId}`, 3600, sessionData); // Expires in 1 hour
```

**Note**: Currently not widely used in the codebase. Consider adding for:
- Session data (1 hour TTL)
- Password reset tokens (15 minutes TTL)
- Temporary data (varies)

---

### 4. Cache-Aside Invalidation

**When**: Read operation, cache miss
**Action**: Fetch from DB, populate cache

```typescript
// Get user from cache or DB
async getUserById(userId: string): Promise<IUserDocument> {
  // 1. Check cache
  let user = await userCache.getUserFromCache(userId);

  // 2. If cache miss, fetch from DB
  if (!user) {
    user = await userService.getUserById(userId);

    // 3. Populate cache
    if (user) {
      await userCache.saveUserToCache(userId, user.uId, user);
    }
  }

  return user;
}
```

---

## Cache Patterns

### Pattern 1: Single Object Cache (HASH)

**Use Case**: Store individual entities (user, post, comment)

```typescript
// Save
await redis.HSET(`entity:${id}`, {
  field1: value1,
  field2: value2,
  field3: value3
});

// Get
const entity = await redis.HGETALL(`entity:${id}`);

// Update single field
await redis.HSET(`entity:${id}`, 'field1', newValue);

// Delete
await redis.DEL(`entity:${id}`);
```

---

### Pattern 2: Sorted Collection Cache (ZSET)

**Use Case**: Store ordered lists (posts by timestamp, users by ID)

```typescript
// Add to sorted set
await redis.ZADD('collection', { score: timestamp, value: `item:${id}` });

// Get range (latest first)
const items = await redis.ZRANGE('collection', 0, 9, { REV: true });

// Get count
const count = await redis.ZCARD('collection');

// Remove item
await redis.ZREM('collection', `item:${id}`);
```

---

### Pattern 3: Sequential Collection Cache (LIST)

**Use Case**: Store ordered sequences (messages, activity logs)

```typescript
// Append to list
await redis.RPUSH(`list:${id}`, JSON.stringify(item));

// Get range
const items = await redis.LRANGE(`list:${id}`, 0, 49); // First 50

// Get count
const count = await redis.LLEN(`list:${id}`);

// Trim (keep only recent)
await redis.LTRIM(`list:${id}`, 0, 99); // Keep 100 items
```

---

### Pattern 4: Counter Cache (STRING)

**Use Case**: Store numeric values (counts, scores)

```typescript
// Set counter
await redis.SET(`counter:${id}`, '0');

// Increment
await redis.INCR(`counter:${id}`);

// Decrement
await redis.DECR(`counter:${id}`);

// Get value
const count = await redis.GET(`counter:${id}`);
```

---

## Performance Metrics

### Cache Hit Rates

**Target Metrics**:
- User profiles: 80-90% hit rate
- Posts feed: 70-85% hit rate
- Messages: 60-75% hit rate
- Reactions: 50-70% hit rate

### Response Time Improvements

**Without Cache**:
- User profile: 300-500ms (MongoDB query)
- Posts feed: 500-1000ms (MongoDB aggregation)
- Messages: 200-400ms (MongoDB query)

**With Cache**:
- User profile: 5-10ms (Redis HGETALL)
- Posts feed: 10-20ms (Redis ZRANGE + HGETALL)
- Messages: 8-15ms (Redis LRANGE)

**Improvement**: 30-100x faster reads

---

## Cache Connection Management

### Redis Connection

```typescript
class BaseCache {
  client: RedisClient;

  constructor(cacheName: string) {
    this.client = createClient({ url: config.REDIS_HOST });
    this.client.on('error', (error: unknown) => {
      log.error(error);
    });
  }

  // Ensure connection is open before operations
  async connect(): Promise<void> {
    if (!this.client.isOpen) {
      await this.client.connect();
    }
  }

  // Close connection gracefully
  async disconnect(): Promise<void> {
    await this.client.quit();
  }
}
```

---

## Best Practices

### 1. Always Check Connection State

```typescript
if (!this.client.isOpen) {
  await this.client.connect();
}
```

### 2. Handle Cache Errors Gracefully

```typescript
try {
  const data = await cache.getData(key);
  return data;
} catch (error) {
  log.error('Cache error:', error);
  // Fallback to database
  return await service.getDataFromDB(key);
}
```

### 3. Serialize Complex Objects

```typescript
// Store
await redis.HSET(`key:${id}`, 'complexField', JSON.stringify(complexObject));

// Retrieve
const rawValue = await redis.HGET(`key:${id}`, 'complexField');
const complexObject = JSON.parse(rawValue);
```

### 4. Use Consistent Key Naming

```
Pattern: entity:identifier:field
Examples:
  - users:507f1f77bcf86cd799439011
  - posts:507f1f77bcf86cd799439012
  - messages:conversation-123
  - followers:count:507f1f77bcf86cd799439011
```

### 5. Batch Operations When Possible

```typescript
// Bad: Multiple round trips
for (const postId of postIds) {
  const post = await redis.HGETALL(`posts:${postId}`);
}

// Good: Pipeline
const pipeline = redis.pipeline();
for (const postId of postIds) {
  pipeline.HGETALL(`posts:${postId}`);
}
const results = await pipeline.exec();
```

---

## Monitoring & Debugging

### Cache Statistics

**Track**:
- Cache hit/miss ratio
- Response times (cache vs DB)
- Memory usage
- Connection errors

**Tools**:
- Redis CLI: `redis-cli INFO stats`
- Redis Commander: GUI for browsing cache
- Application logs: Track cache operations

### Common Issues

**Issue 1: Cache not updating**
- Check queue workers are running
- Verify Redis connection is open
- Check for errors in worker logs

**Issue 2: Stale data in cache**
- Implement cache invalidation on writes
- Consider TTL for infrequently changing data
- Manual cache clear if needed: `redis-cli FLUSHDB`

**Issue 3: Memory issues**
- Monitor Redis memory usage: `redis-cli INFO memory`
- Implement TTL on appropriate keys
- Use LTRIM on lists to limit size
- Consider eviction policies (LRU)

---

## Future Improvements

1. **Add TTL to more cache keys** - Auto-expire stale data
2. **Implement cache warming** - Pre-populate cache on startup
3. **Add cache metrics** - Track hit/miss rates programmatically
4. **Use Redis Cluster** - Horizontal scaling for high traffic
5. **Implement cache versioning** - Handle schema changes gracefully
6. **Add cache compression** - Reduce memory usage for large objects
