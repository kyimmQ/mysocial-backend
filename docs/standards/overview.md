# MySocial Backend - Code Standards

## Overview

This document serves as an index to the coding standards, conventions, and best practices used throughout the MySocial Backend codebase. All contributors should follow these guidelines to maintain code quality, consistency, and maintainability.

---

## Standards Documentation

Our code standards are organized into focused documents:

### 1. [TypeScript Standards](./typescript.md)
- Type annotations and interface naming
- Type vs Interface usage
- Strict mode configuration
- Path aliases
- Naming conventions (files, classes, variables, functions)

### 2. [Architecture Patterns](./architecture-patterns.md)
- Controller pattern
- Service pattern
- Cache pattern
- Queue pattern
- Worker pattern
- Error handling
- Async/await standards
- Validation standards

### 3. [Testing Standards](./testing.md)
- Test file structure
- Test coverage requirements
- Mock usage
- Code quality tools (ESLint, Prettier)
- Logging standards (Bunyan)

### 4. [Security & Performance](./security-performance.md)
- Password hashing
- JWT tokens
- Input sanitization
- CORS configuration
- Caching strategies
- Query optimization
- Git commit standards

---

## Quick Reference

### Most Important Rules

1. **TypeScript**: Always use explicit types for function parameters and return types
2. **Naming**: Use PascalCase for classes, camelCase for variables/functions, kebab-case for files
3. **Async**: Always use async/await, never callbacks
4. **Validation**: Use @joiValidation decorator on all controller methods
5. **Caching**: Cache-first for reads, write-through for writes
6. **Testing**: Minimum 80% code coverage for all controllers
7. **Errors**: Use custom error classes, never throw raw errors
8. **Logging**: Use Bunyan logger with structured logging

### Common Patterns

**Controller Flow**:
```
Request → Validation → Cache Check → Business Logic → Cache Write → Queue Job → Socket Emit → Response
```

**Service Pattern**:
```typescript
// Singleton export
export const userService: UserService = new UserService();
```

**Cache Pattern**:
```typescript
// Extends BaseCache
export class UserCache extends BaseCache {
  constructor() {
    super('userCache');
  }
}
```

**Queue Pattern**:
```typescript
// Extends BaseQueue
class UserQueue extends BaseQueue {
  constructor() {
    super('users');
    this.processJob('addUserToDB', 5, userWorker.addUserToDB);
  }
}
```

---

## Enforcement

These standards are enforced through:

- **ESLint** - Linting rules (`.eslintrc.json`)
- **Prettier** - Code formatting (`.prettierrc.json`)
- **TypeScript** - Type checking (`tsconfig.json`)
- **Jest** - Test coverage thresholds (`jest.config.ts`)
- **Code Review** - Manual review process

---

## Tools

### Check Code Quality
```bash
npm run lint:check        # Run ESLint
npm run prettier:check    # Check formatting
npm run test              # Run tests with coverage
```

### Auto-Fix Issues
```bash
npm run lint:fix          # Auto-fix linting errors
npm run prettier:fix      # Auto-format code
```

---

## Getting Help

- Review the specific standards document for detailed guidance
- Check existing code for examples
- Ask in team discussions for clarification
- Refer to the [Codebase Summary](../codebase-summary.md) for architecture context

---

## Contributing

When adding new patterns or standards:

1. Document them in the appropriate standards file
2. Provide code examples
3. Update this overview with cross-references
4. Ensure consistency with existing patterns
5. Get team review before enforcing

---

These code standards ensure consistency, maintainability, and quality across the MySocial Backend codebase. All contributors should familiarize themselves with these guidelines and apply them consistently in their work.
