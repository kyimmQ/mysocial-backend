# Architecture Patterns

## Controller Pattern

All controllers follow this structure:

```typescript
import { Request, Response } from 'express';
import HTTP_STATUS from 'http-status-codes';
import { joiValidation } from '@global/decorators/joi-validation.decorators';
import { signupSchema } from '@auth/schemes/signup';

export class SignUp {
  @joiValidation(signupSchema)
  public async create(req: Request, res: Response): Promise<void> {
    // 1. Extract and validate data from request
    const { username, email, password } = req.body;

    // 2. Check cache (if applicable)
    const existingUser = await userCache.getUserFromCache(username);

    // 3. Perform business logic
    const authData = SignUp.prototype.signupData(data);

    // 4. Write to cache
    await userCache.saveUserToCache(userId, userDataForCache);

    // 5. Add job to queue for async DB write
    await authQueue.addAuthUserJob('addAuthUserToDB', { value: authData });

    // 6. Emit real-time event (if applicable)
    socketIOPostObject.emit('add post', postData);

    // 7. Return response
    res.status(HTTP_STATUS.CREATED).json({
      message: 'User created successfully',
      user: userData,
      token: userJwt
    });
  }

  private signupData(data: ISignUpData): IAuthDocument {
    // Helper methods as private
  }
}
```

### Controller Guidelines

- One class per controller file
- Public methods for route handlers
- Private methods for helper logic
- Use `@joiValidation` decorator for validation
- Return HTTP status codes from `http-status-codes`
- Never call database services directly - use cache first
- Always queue async operations with Bull
- Emit Socket.IO events for real-time updates

---

## Service Pattern

All database services follow this structure:

```typescript
import { UserModel } from '@user/models/user.schema';
import { IUserDocument } from '@user/interfaces/user.interface';

class UserService {
  public async addUserData(data: IUserDocument): Promise<void> {
    await UserModel.create(data);
  }

  public async getUserById(userId: string): Promise<IUserDocument> {
    const user: IUserDocument = await UserModel.findById(userId).exec();
    return user;
  }

  public async getUserByUsername(username: string): Promise<IUserDocument> {
    const user: IUserDocument = await UserModel.findOne({ username }).exec();
    return user;
  }

  public async updateUser(userId: string, data: Partial<IUserDocument>): Promise<void> {
    await UserModel.updateOne({ _id: userId }, { $set: data }).exec();
  }

  public async deleteUser(userId: string): Promise<void> {
    await UserModel.deleteOne({ _id: userId }).exec();
  }
}

export const userService: UserService = new UserService();
```

### Service Guidelines

- One service class per domain (User, Post, Auth, etc.)
- Export as singleton instance
- Only database operations - no business logic
- Use Mongoose models for queries
- Always call `.exec()` on queries
- Return typed results
- Handle errors at service level

---

## Cache Pattern

All cache services follow this structure:

```typescript
import { BaseCache } from '@service/redis/base.cache';
import { IUserDocument } from '@user/interfaces/user.interface';

export class UserCache extends BaseCache {
  constructor() {
    super('userCache');
  }

  public async saveUserToCache(key: string, userUId: string, createdUser: IUserDocument): Promise<void> {
    const createdAt = new Date();
    const {
      _id,
      uId,
      username,
      email,
      avatarColor,
      // ... other fields
    } = createdUser;

    const dataToSave = {
      '_id': `${_id}`,
      'uId': `${uId}`,
      'username': `${username}`,
      'email': `${email}`,
      'avatarColor': `${avatarColor}`,
      'createdAt': `${createdAt}`,
      // ... other fields
    };

    try {
      if (!this.client.isOpen) {
        await this.client.connect();
      }
      await this.client.ZADD('user', { score: parseInt(userUId, 10), value: `${key}` });
      await this.client.HSET(`users:${key}`, dataToSave);
    } catch (error) {
      log.error(error);
      throw error;
    }
  }

  public async getUserFromCache(userId: string): Promise<IUserDocument | null> {
    try {
      if (!this.client.isOpen) {
        await this.client.connect();
      }

      const response: IUserDocument = await this.client.HGETALL(`users:${userId}`) as unknown as IUserDocument;

      response.createdAt = new Date(Helpers.parseJson(`${response.createdAt}`));
      response.postsCount = Helpers.parseJson(`${response.postsCount}`);
      // ... parse other fields

      return response;
    } catch (error) {
      log.error(error);
      return null;
    }
  }
}
```

### Cache Guidelines

- Extend `BaseCache` class
- Name cache in constructor (for logging)
- Convert all values to strings for Redis
- Parse JSON strings back to objects
- Check connection before operations
- Use appropriate Redis data structures:
  - HASH for objects
  - ZSET for sorted lists
  - LIST for collections
  - STRING for simple values
- Always handle errors gracefully
- Return null on cache miss

---

## Queue Pattern

All queues follow this structure:

```typescript
import { BaseQueue } from '@service/queues/base.queue';
import { IAuthJob } from '@auth/interfaces/auth.interface';

class AuthQueue extends BaseQueue {
  constructor() {
    super('auth');
    this.processJob('addAuthUserToDB', 5, authWorker.addAuthUserToDB);
  }

  public addAuthUserJob(name: string, data: IAuthJob): void {
    this.addJob(name, data);
  }
}

export const authQueue: AuthQueue = new AuthQueue();
```

### Queue Guidelines

- Extend `BaseQueue` class
- Name queue in constructor
- Register workers in constructor with `processJob()`
- Set concurrency (typically 5)
- Create typed job methods
- Export as singleton
- Job names should be descriptive actions (e.g., 'addUserToDB', 'sendEmail')

---

## Worker Pattern

All workers follow this structure:

```typescript
import { DoneCallback, Job } from 'bull';
import Logger from 'bunyan';
import { config } from '@root/config';
import { authService } from '@service/db/auth.service';

const log: Logger = config.createLogger('authWorker');

class AuthWorker {
  async addAuthUserToDB(job: Job, done: DoneCallback): Promise<void> {
    try {
      const { value } = job.data;
      await authService.createAuthUser(value);
      job.progress(100);
      done(null, job.data);
    } catch (error) {
      log.error(error);
      done(error as Error);
    }
  }
}

export const authWorker: AuthWorker = new AuthWorker();
```

### Worker Guidelines

- One worker class per queue
- Create Bunyan logger instance
- Extract data from `job.data`
- Call appropriate service method
- Update job progress to 100% on success
- Call `done(null, result)` on success
- Call `done(error)` on failure
- Log all errors
- Export as singleton

---

## Error Handling

### Custom Error Classes

All errors extend the base `CustomError` class:

```typescript
export abstract class CustomError extends Error {
  abstract statusCode: number;
  abstract status: string;

  constructor(message: string) {
    super(message);
  }

  serializeErrors(): IError {
    return {
      message: this.message,
      statusCode: this.statusCode,
      status: this.status
    };
  }
}
```

**Available error classes**:
- `BadRequestError` (400) - Invalid input
- `NotFoundError` (404) - Resource not found
- `NotAuthorizedError` (401) - Authentication failed
- `FileTooLargeError` (413) - File size exceeded
- `ServerError` (500) - Internal server error

**Usage**:
```typescript
import { BadRequestError } from '@global/helpers/error-handler';

if (!username) {
  throw new BadRequestError('Username is required');
}
```

### Try-Catch Blocks

**Use try-catch for**:
- Database operations
- Redis operations
- External API calls
- File operations

**Example**:
```typescript
try {
  const user = await UserModel.findById(userId);
  if (!user) {
    throw new NotFoundError('User not found');
  }
  return user;
} catch (error) {
  log.error(error);
  throw error;
}
```

### Express Async Errors

The application uses `express-async-errors` to automatically catch async errors in route handlers:

```typescript
import 'express-async-errors';

// No need for try-catch in route handlers
app.get('/users/:id', async (req, res) => {
  const user = await userService.getUserById(req.params.id); // Errors auto-caught
  res.json(user);
});
```

---

## Async/Await Standards

### Always Use Async/Await

**Never use callbacks or raw promises**:

```typescript
// Good
async function getUser(userId: string): Promise<IUserDocument> {
  const user = await UserModel.findById(userId);
  return user;
}

// Bad
function getUser(userId: string, callback: Function): void {
  UserModel.findById(userId, (err, user) => {
    callback(err, user);
  });
}
```

### Promise.all for Parallel Operations

```typescript
// Good - Parallel execution
const [user, posts, followers] = await Promise.all([
  userService.getUserById(userId),
  postService.getPostsByUserId(userId),
  followerService.getFollowersByUserId(userId)
]);

// Bad - Sequential execution
const user = await userService.getUserById(userId);
const posts = await postService.getPostsByUserId(userId);
const followers = await followerService.getFollowersByUserId(userId);
```

### Error Handling in Async Functions

```typescript
async function createPost(data: IPostDocument): Promise<void> {
  try {
    await postService.createPost(data);
    await postCache.savePostToCache(data);
  } catch (error) {
    log.error('Error creating post:', error);
    throw new ServerError('Failed to create post');
  }
}
```

---

## Validation Standards

### Joi Schemas

All request validation uses Joi schemas:

```typescript
import Joi, { ObjectSchema } from 'joi';

const signupSchema: ObjectSchema = Joi.object().keys({
  username: Joi.string().min(4).max(8).required().messages({
    'string.base': 'Username must be of type string',
    'string.min': 'Username must be at least 4 characters',
    'string.max': 'Username must be at most 8 characters',
    'string.empty': 'Username is required'
  }),
  password: Joi.string().min(4).max(8).required().messages({
    'string.base': 'Password must be of type string',
    'string.min': 'Password must be at least 4 characters',
    'string.max': 'Password must be at most 8 characters',
    'string.empty': 'Password is required'
  }),
  email: Joi.string().email().required().messages({
    'string.base': 'Email must be of type string',
    'string.email': 'Email must be valid',
    'string.empty': 'Email is required'
  }),
  avatarColor: Joi.string().optional(),
  avatarImage: Joi.string().optional()
});
```

### Validation Decorator

Use the `@joiValidation` decorator for controller methods:

```typescript
import { joiValidation } from '@global/decorators/joi-validation.decorators';
import { signupSchema } from '@auth/schemes/signup';

export class SignUp {
  @joiValidation(signupSchema)
  public async create(req: Request, res: Response): Promise<void> {
    // Request body is already validated
    const { username, email, password } = req.body;
  }
}
```

---

## Summary

These architecture patterns provide:
- Consistent structure across the codebase
- Clear separation of concerns (Controller → Cache → Queue → Worker → Service → DB)
- Predictable data flow
- Easy testing through well-defined interfaces
- Scalability through async processing
- Type safety with TypeScript

All code should follow these patterns to maintain architectural consistency.
