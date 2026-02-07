# TypeScript Standards

## Type Annotations

### Always use explicit types for

- Function parameters
- Function return types
- Class properties
- Exported constants

**Example**:
```typescript
// Good
public async createPost(data: IPostDocument): Promise<void> {
  // implementation
}

// Bad
public async createPost(data) {
  // implementation
}
```

---

## Interface Naming

- Use `I` prefix for interfaces: `IUserDocument`, `IPostJob`, `IAuthPayload`
- Use descriptive names that indicate purpose: `ISignUpData`, `IReactionJob`

**Example**:
```typescript
export interface IUserDocument extends Document {
  _id: string | ObjectId;
  username: string;
  email: string;
  password?: string;
  avatarColor: string;
}
```

---

## Type vs Interface

### Use `interface` for

- Object shapes
- API contracts
- Extending other types

### Use `type` for

- Union types
- Intersection types
- Mapped types
- Function types

**Example**:
```typescript
// Interface for object shapes
interface IUser {
  username: string;
  email: string;
}

// Type for unions
type ReactionType = 'like' | 'love' | 'happy' | 'wow' | 'sad' | 'angry';

// Type for functions
type EmailCallback = (data: IEmailJob) => Promise<void>;
```

---

## Strict Mode

All TypeScript files are compiled with `strict: true` in `tsconfig.json`:
- `noImplicitAny`: true
- `strictNullChecks`: true
- `strictFunctionTypes`: true
- `strictPropertyInitialization`: true
- `noImplicitThis`: true
- `alwaysStrict`: true

---

## Path Aliases

Use TypeScript path aliases for imports (configured in `tsconfig.json`):

```typescript
// Good
import { UserCache } from '@service/redis/user.cache';
import { authMiddleware } from '@global/helpers/auth-middleware';
import { IAuthDocument } from '@auth/interfaces/auth.interface';

// Bad
import { UserCache } from '../../../shared/services/redis/user.cache';
```

### Available aliases

- `@auth/*` - Auth feature
- `@user/*` - User feature
- `@post/*` - Post feature
- `@reaction/*` - Reaction feature
- `@comment/*` - Comment feature
- `@follower/*` - Follower feature
- `@chat/*` - Chat feature
- `@notification/*` - Notification feature
- `@image/*` - Image feature
- `@global/*` - Global helpers
- `@service/*` - Shared services
- `@socket/*` - Socket.IO handlers
- `@worker/*` - Queue workers
- `@root/*` - Root directory

---

## Naming Conventions

### Files

**Controllers**: Kebab-case with action prefix
- `signup.ts`, `signin.ts`, `create-post.ts`, `get-profile.ts`

**Models**: Kebab-case with `.schema` suffix
- `user.schema.ts`, `post.schema.ts`, `auth.schema.ts`

**Services**: Kebab-case with `.service` suffix
- `user.service.ts`, `post.service.ts`, `auth.service.ts`

**Caches**: Kebab-case with `.cache` suffix
- `user.cache.ts`, `post.cache.ts`, `message.cache.ts`

**Queues**: Kebab-case with `.queue` suffix
- `user.queue.ts`, `post.queue.ts`, `email.queue.ts`

**Workers**: Kebab-case with `.worker` suffix
- `user.worker.ts`, `post.worker.ts`, `email.worker.ts`

**Routes**: CamelCase with `Routes` suffix
- `authRoutes.ts`, `postRoutes.ts`, `userRoutes.ts`

**Interfaces**: Kebab-case with `.interface` suffix
- `user.interface.ts`, `post.interface.ts`, `auth.interface.ts`

**Schemas (Joi)**: Kebab-case or singular
- `signup.ts`, `signin.ts`, `post.schemes.ts`, `comment.ts`

**Tests**: Same as source file with `.test` suffix
- `signup.test.ts`, `create-post.test.ts`, `get-profile.test.ts`

**Mocks**: Singular kebab-case with `.mock` suffix
- `user.mock.ts`, `post.mock.ts`, `auth.mock.ts`

---

### Classes

**PascalCase** for all classes:
- `SignUp`, `CreatePost`, `UserCache`, `PostService`, `AuthQueue`

**Controller classes**: Action + feature name
- `SignUp`, `SignIn`, `CreatePost`, `GetProfile`, `AddComment`

**Service classes**: Feature + `Service`
- `UserService`, `PostService`, `AuthService`, `ChatService`

**Cache classes**: Feature + `Cache`
- `UserCache`, `PostCache`, `MessageCache`, `ReactionCache`

**Queue classes**: Feature + `Queue`
- `UserQueue`, `PostQueue`, `EmailQueue`, `NotificationQueue`

**Worker classes**: Feature + `Worker`
- `UserWorker`, `PostWorker`, `EmailWorker`, `AuthWorker`

---

### Variables

**camelCase** for variables and function parameters:
- `userId`, `postData`, `avatarColor`, `createdAt`

**UPPER_SNAKE_CASE** for constants:
- `DEFAULT_PAGE_SIZE`, `MAX_FILE_SIZE`, `JWT_EXPIRATION`

**Example**:
```typescript
const DEFAULT_PAGE_SIZE = 10;
const MAX_RETRY_ATTEMPTS = 3;

let userId: string;
let postData: IPostDocument;
```

---

### Functions

**camelCase** for function names:
- `getUserById()`, `createPost()`, `sendEmail()`, `validateToken()`

**Verb + noun pattern** for clarity:
- `getUser()`, `createPost()`, `deleteComment()`, `updateProfile()`

**Boolean functions**: Use `is`, `has`, `can`, `should` prefixes
- `isValidEmail()`, `hasPermission()`, `canDeletePost()`, `shouldSendNotification()`

---

### Mongoose Models

**PascalCase** for model names (singular):
- `AuthModel`, `UserModel`, `PostModel`, `CommentModel`

**Model variable naming**:
```typescript
const AuthModel: Model<IAuthDocument> = model<IAuthDocument>('Auth', authSchema, 'Auth');
```

---

## Summary

These TypeScript standards ensure:
- Type safety across the codebase
- Consistent naming conventions
- Better IDE autocomplete and IntelliSense
- Easier code navigation
- Reduced runtime errors
- Self-documenting code

All code must follow these conventions to maintain quality and consistency.
