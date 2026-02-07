# Testing Standards

## Test File Structure

```typescript
import { Request, Response } from 'express';
import { SignUp } from '@auth/controllers/signup';
import { authMock, authMockRequest, authMockResponse } from '@root/mocks/auth.mock';
import { CustomError } from '@global/helpers/error-handler';

jest.useFakeTimers();

describe('SignUp', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  afterEach(() => {
    jest.clearAllMocks();
    jest.clearAllTimers();
  });

  it('should throw an error if username is not available', () => {
    const req: Request = authMockRequest({}, { username: '', email: 'test@test.com', password: 'qwerty' }) as Request;
    const res: Response = authMockResponse();

    SignUp.prototype.create(req, res).catch((error: CustomError) => {
      expect(error.statusCode).toEqual(400);
      expect(error.serializeErrors().message).toEqual('Username is a required field');
    });
  });

  it('should create a new user and return success', async () => {
    const req: Request = authMockRequest({}, authMock) as Request;
    const res: Response = authMockResponse();

    await SignUp.prototype.create(req, res);

    expect(res.status).toHaveBeenCalledWith(201);
    expect(res.json).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'User created successfully',
        user: expect.any(Object),
        token: expect.any(String)
      })
    );
  });
});
```

---

## Test Organization

### File Location

Tests are co-located with source code in `test/` subdirectories:

```
src/features/auth/controllers/
  ├── signup.ts
  └── test/
      └── signup.test.ts
```

### Test Naming

- Use descriptive test names starting with "should"
- Group related tests with `describe` blocks
- One test file per controller/service/utility

**Example**:
```typescript
describe('UserService', () => {
  describe('getUserById', () => {
    it('should return user when found', async () => { /* ... */ });
    it('should throw NotFoundError when user does not exist', async () => { /* ... */ });
  });
});
```

---

## Test Coverage

### Requirements

- **Minimum 80% code coverage** for all controllers
- Test both happy paths and error cases
- Cover edge cases and boundary conditions

### Coverage Commands

```bash
npm run test              # Run tests with coverage
```

### What to Test

**Controllers**:
- ✅ Validation errors
- ✅ Successful operations
- ✅ Cache interactions
- ✅ Queue job additions
- ✅ Socket.IO emissions
- ✅ HTTP status codes
- ✅ Response structure

**Services**:
- ✅ Database queries
- ✅ Null/undefined handling
- ✅ Error conditions

**Utilities**:
- ✅ Input/output transformations
- ✅ Edge cases
- ✅ Error handling

---

## Mocking

### Mock Structure

All mocks are centralized in `/src/mocks/`:

```typescript
// user.mock.ts
export const existingUser = {
  _id: '60263f14648fed5246e322d9',
  username: 'Manny',
  email: 'manny@test.com',
  avatarColor: '#9c27b0',
  uId: '1621613119252066',
  // ... other fields
};

export const userMockRequest = (sessionData: IJWT, body: IBody, currentUser?: AuthPayload | null, params?: IParams) => ({
  session: sessionData,
  body,
  params,
  currentUser
});

export const userMockResponse = (): Response => {
  const res: Response = {} as Response;
  res.status = jest.fn().mockReturnValue(res);
  res.json = jest.fn().mockReturnValue(res);
  return res;
};
```

### Mock Usage

```typescript
import { authMockRequest, authMockResponse } from '@root/mocks/auth.mock';

const req: Request = authMockRequest({}, { username: 'test', email: 'test@test.com' }) as Request;
const res: Response = authMockResponse();

await SignUp.prototype.create(req, res);

expect(res.status).toHaveBeenCalledWith(201);
```

### Mock External Dependencies

**Database**:
```typescript
jest.mock('@service/db/user.service');
```

**Redis**:
```typescript
jest.mock('@service/redis/user.cache');
```

**Queues**:
```typescript
jest.mock('@service/queues/user.queue');
```

**Socket.IO**:
```typescript
jest.mock('@socket/user');
```

---

## Code Quality Tools

### ESLint

Configuration in `.eslintrc.json`:
- TypeScript ESLint parser
- Recommended rules
- No unused variables
- No console.log in production
- Consistent return statements

**Run linting**:
```bash
npm run lint:check    # Check for errors
npm run lint:fix      # Auto-fix errors
```

### Prettier

Configuration in `.prettierrc.json`:
- 2-space indentation
- Single quotes
- Semicolons required
- Trailing commas: ES5
- Print width: 120

**Run formatting**:
```bash
npm run prettier:check  # Check formatting
npm run prettier:fix    # Auto-format
```

### Editor Config

Configuration in `.editorconfig`:
- UTF-8 charset
- LF line endings
- Trim trailing whitespace
- Insert final newline

---

## Logging Standards

### Bunyan Logger

Use Bunyan for structured logging:

```typescript
import Logger from 'bunyan';
import { config } from '@root/config';

const log: Logger = config.createLogger('userService');

log.info('User created successfully', { userId, username });
log.error('Failed to create user', { error, userId });
log.warn('User not found in cache', { userId });
log.debug('Cache hit', { key, value });
```

### Log Levels

- `fatal` - Application crash
- `error` - Errors that need attention
- `warn` - Warnings (degraded functionality)
- `info` - Important information
- `debug` - Debugging information
- `trace` - Detailed tracing

### What to Log

**Do log**:
- User actions (signup, login, post creation)
- Errors with context
- Cache hits/misses
- Queue job processing
- API requests/responses

**Don't log**:
- Passwords or sensitive data
- Full request bodies (unless debugging)
- Excessive debug info in production

### Production vs Development

**Development**:
- Log level: `debug`
- Output: Pretty-printed to console

**Production**:
- Log level: `info`
- Output: JSON format piped through Bunyan
- Centralized logging (optional)

---

## Test Best Practices

### 1. Arrange, Act, Assert (AAA)

```typescript
it('should create a user', async () => {
  // Arrange
  const req = authMockRequest({}, { username: 'test' });
  const res = authMockResponse();

  // Act
  await SignUp.prototype.create(req, res);

  // Assert
  expect(res.status).toHaveBeenCalledWith(201);
});
```

### 2. One Assertion Per Test

```typescript
// Good
it('should return 201 status', async () => {
  await SignUp.prototype.create(req, res);
  expect(res.status).toHaveBeenCalledWith(201);
});

it('should return user data', async () => {
  await SignUp.prototype.create(req, res);
  expect(res.json).toHaveBeenCalledWith(expect.objectContaining({ user: expect.any(Object) }));
});

// Bad
it('should create user and return response', async () => {
  await SignUp.prototype.create(req, res);
  expect(res.status).toHaveBeenCalledWith(201);
  expect(res.json).toHaveBeenCalledWith(expect.objectContaining({ user: expect.any(Object) }));
  expect(res.json).toHaveBeenCalledWith(expect.objectContaining({ token: expect.any(String) }));
});
```

### 3. Clear Test Names

```typescript
// Good
it('should throw BadRequestError when username is empty', async () => { /* ... */ });

// Bad
it('test signup', async () => { /* ... */ });
```

### 4. Use beforeEach for Setup

```typescript
describe('UserService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // tests...
});
```

### 5. Test Edge Cases

```typescript
describe('validateEmail', () => {
  it('should accept valid email', () => { /* ... */ });
  it('should reject email without @', () => { /* ... */ });
  it('should reject email without domain', () => { /* ... */ });
  it('should reject empty email', () => { /* ... */ });
  it('should reject null email', () => { /* ... */ });
});
```

---

## Continuous Integration

Tests run automatically in CircleCI pipeline:

```yaml
- run:
    name: Run tests
    command: npm run test

- run:
    name: Upload coverage to Codecov
    command: npx codecov
```

### Coverage Reporting

- Local: Terminal output + HTML report in `/coverage`
- CI/CD: Codecov integration with badges

---

## Summary

These testing standards ensure:
- High code quality
- Reliable functionality
- Easy refactoring
- Consistent test structure
- Comprehensive coverage
- Early bug detection

All code must include tests before merging to maintain quality standards.
