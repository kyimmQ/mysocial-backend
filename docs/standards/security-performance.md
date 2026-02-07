# Security & Performance Standards

## Security Standards

### Password Hashing

Always use bcrypt with 10 salt rounds:

```typescript
import { hash } from 'bcryptjs';

const hashedPassword: string = await hash(password, 10);
```

**In Mongoose schemas**:
```typescript
authSchema.pre('save', async function(this: IAuthDocument) {
  if (this.isModified('password')) {
    const hashedPassword: string = await hash(this.password as string, 10);
    this.password = hashedPassword;
  }
});
```

**Password comparison**:
```typescript
authSchema.methods.comparePassword = async function(password: string): Promise<boolean> {
  const hashedPassword: string = (this as IAuthDocument).password!;
  return compare(password, hashedPassword);
};
```

---

### JWT Tokens

Generate JWT tokens with expiration:

```typescript
import JWT from 'jsonwebtoken';
import { config } from '@root/config';

const token: string = JWT.sign(
  { userId, username, email },
  config.JWT_TOKEN!,
  { expiresIn: '1d' }
);
```

**Token verification**:
```typescript
const payload: AuthPayload = JWT.verify(token, config.JWT_TOKEN!) as AuthPayload;
```

**Session storage**:
```typescript
req.session = { jwt: token };
```

---

### Input Sanitization

All inputs validated with Joi before processing:

```typescript
@joiValidation(schema)
public async method(req: Request, res: Response): Promise<void> {
  // Input is already validated and sanitized
}
```

**Regex sanitization** for search queries:
```typescript
import { Helpers } from '@global/helpers/helpers';

const sanitizedQuery = Helpers.escapeRegex(searchQuery);
const users = await UserModel.find({
  username: { $regex: sanitizedQuery, $options: 'i' }
});
```

---

### CORS Configuration

Restrict CORS to trusted origins:

```typescript
const corsOptions = {
  origin: config.CLIENT_URL,
  credentials: true,
  optionsSuccessStatus: 200
};

app.use(cors(corsOptions));
```

---

### Security Middleware

**Helmet** - Security headers:
```typescript
app.use(helmet());
```

**HPP** - HTTP Parameter Pollution protection:
```typescript
app.use(hpp());
```

**Cookie Session** - Encrypted sessions:
```typescript
app.use(cookieSession({
  name: 'session',
  keys: [config.SECRET_KEY_ONE!, config.SECRET_KEY_TWO!],
  maxAge: 7 * 24 * 60 * 60 * 1000, // 7 days
  secure: config.NODE_ENV === 'production',
  sameSite: 'none'
}));
```

---

### Sensitive Data Handling

**Never log**:
- Passwords
- JWT tokens (except for debugging)
- API keys
- Personal identifiable information (PII)

**Exclude from responses**:
```typescript
authSchema.set('toJSON', {
  transform(_doc, ret) {
    delete ret.password;
    return ret;
  }
});
```

---

## Performance Standards

### Caching Strategy

#### Cache-First for Reads

```typescript
// 1. Check cache
let user = await userCache.getUserFromCache(userId);

// 2. If not in cache, fetch from DB
if (!user) {
  user = await userService.getUserById(userId);

  // 3. Populate cache
  await userCache.saveUserToCache(userId, user);
}

return user;
```

#### Write-Through for Writes

```typescript
// 1. Write to cache immediately
await postCache.savePostToCache(postId, postData);

// 2. Queue job for DB write
await postQueue.addPostJob('savePostToDB', { key: postId, value: postData });

// 3. Return success (don't wait for DB)
res.status(201).json({ message: 'Post created' });
```

---

### Query Optimization

#### Use lean() for Read-Only Queries

```typescript
// Returns plain JavaScript objects instead of Mongoose documents
const posts = await PostModel.find({ userId }).lean().exec();
```

#### Use select() to Limit Fields

```typescript
const user = await UserModel.findById(userId).select('username email avatarColor').exec();
```

#### Use Indexes for Frequently Queried Fields

```typescript
// In schema definition
userSchema.index({ username: 1 });
userSchema.index({ email: 1 });
postSchema.index({ userId: 1, createdAt: -1 });
```

#### Pagination for Large Datasets

```typescript
const PAGE_SIZE = 10;
const skip = (page - 1) * PAGE_SIZE;

const posts = await PostModel
  .find()
  .sort({ createdAt: -1 })
  .skip(skip)
  .limit(PAGE_SIZE)
  .exec();
```

#### Aggregation for Complex Queries

```typescript
const result = await UserModel.aggregate([
  { $match: { _id: new mongoose.Types.ObjectId(userId) } },
  {
    $lookup: {
      from: 'Auth',
      localField: 'authId',
      foreignField: '_id',
      as: 'authData'
    }
  },
  { $unwind: '$authData' },
  { $project: { password: 0 } }
]);
```

---

### Redis Optimization

#### Use Appropriate Data Structures

**HASH** for objects:
```typescript
await client.HSET(`users:${userId}`, userData);
```

**ZSET** for sorted lists:
```typescript
await client.ZADD('posts', { score: parseInt(userId), value: postId });
```

**LIST** for collections:
```typescript
await client.LPUSH(`comments:${postId}`, JSON.stringify(comment));
```

#### Use Multi Commands for Atomic Operations

```typescript
const multi = client.multi();
multi.ZADD('posts', { score: userId, value: postId });
multi.HSET(`posts:${postId}`, postData);
await multi.exec();
```

#### Set Expiration for Temporary Data

```typescript
await client.SETEX(`reset-token:${token}`, 3600, userId); // 1 hour
```

---

### Socket.IO Optimization

#### Use Rooms for Targeted Broadcasting

```typescript
// Join room
socket.join(`user:${userId}`);

// Emit to specific room
io.to(`user:${userId}`).emit('notification', data);
```

#### Limit Payload Size

```typescript
// Bad - Sending entire user object
io.emit('update user', fullUserObject);

// Good - Sending only necessary fields
io.emit('update user', { userId, username, avatarColor });
```

---

### Queue Optimization

#### Batch Similar Operations

```typescript
// Bad - Multiple individual jobs
for (const user of users) {
  await emailQueue.addEmailJob('sendEmail', { userId: user.id });
}

// Good - Single batch job
await emailQueue.addEmailJob('sendBulkEmail', { userIds: users.map(u => u.id) });
```

#### Set Appropriate Concurrency

```typescript
// High concurrency for I/O operations
this.processJob('sendEmail', 10, emailWorker.sendEmail);

// Low concurrency for CPU-intensive operations
this.processJob('processVideo', 2, videoWorker.processVideo);
```

---

### Image Upload Optimization

#### Compress Before Upload

```typescript
// Use Cloudinary transformations
const result = await cloudinary.uploader.upload(image, {
  quality: 'auto',
  fetch_format: 'auto'
});
```

#### Lazy Load Images

```typescript
// Return CDN URLs, let client lazy load
res.json({
  images: posts.map(p => ({
    url: `https://res.cloudinary.com/${cloudName}/image/upload/${p.imgId}`,
    thumbnail: `https://res.cloudinary.com/${cloudName}/image/upload/w_200,h_200/${p.imgId}`
  }))
});
```

---

### API Response Optimization

#### Use Compression

```typescript
import compression from 'compression';
app.use(compression());
```

#### Limit Response Size

```typescript
// Bad - Returning everything
const posts = await PostModel.find().exec();

// Good - Pagination and field selection
const posts = await PostModel
  .find()
  .select('userId username post createdAt')
  .limit(10)
  .exec();
```

---

## Git Commit Standards

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting)
- `refactor` - Code refactoring
- `test` - Test additions or changes
- `chore` - Build process or tool changes

**Example**:
```
feat(auth): add password reset functionality

- Add forgot password endpoint
- Add reset password endpoint
- Add email template for password reset
- Add tests for password reset flow

Closes #123
```

### Commit Best Practices

- Write clear, concise commit messages
- Use imperative mood ("Add feature" not "Added feature")
- Keep subject line under 72 characters
- Reference issue numbers in footer
- Group related changes in single commit
- Avoid mixing refactoring with features

---

## Monitoring & Debugging

### API Monitoring

**Swagger Stats** dashboard at `/api-monitoring`:
- Request/response times
- Error rates
- Endpoint performance
- HTTP status codes

### Queue Monitoring

**Bull Board** UI at `/queues`:
- Job status (active, completed, failed)
- Queue metrics
- Job retry counts
- Processing times

### Redis Monitoring

**Redis Commander** (external tool):
```bash
npm install -g redis-commander
redis-commander
```

### MongoDB Monitoring

**MongoDB Compass** (desktop app):
- Visual query builder
- Index analysis
- Performance insights
- Document editing

---

## Summary

These security and performance standards ensure:
- **Security**: Protection against common vulnerabilities
- **Performance**: Fast response times and efficient resource usage
- **Scalability**: Ability to handle growth
- **Reliability**: Consistent behavior under load
- **Maintainability**: Easy to debug and optimize

All code must follow these standards to maintain production-grade quality.
