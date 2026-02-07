# MySocial Backend - Node.js Implementation (Learning Branch)

> 🎓 **Learning Project**: This branch is for learning backend development by implementing MySocial in Node.js/TypeScript from scratch.
>
> **Branch Strategy**:
> - `node-master-onprogress` - Node.js/TypeScript implementation (this branch - starting from scratch)
> - `main` - Golang implementation (planned)
> - `node-main` - Complete Node.js reference implementation

---

## Project Overview

MySocial Backend is a production-ready, full-featured real-time social media backend application. This branch is for implementing it from scratch using **Node.js** and **TypeScript** with modern patterns and best practices.

### Learning Goals

This implementation will cover:
- ✅ REST API development with Express.js
- ✅ Real-time WebSocket communication (Socket.IO)
- ✅ MongoDB integration with Mongoose ODM
- ✅ Redis caching strategies
- ✅ Async job processing with Bull queues
- ✅ JWT authentication
- ✅ TypeScript patterns and best practices
- ✅ Testing with Jest
- ✅ Docker containerization and AWS deployment

---

## Technology Stack

### Core Technologies
- **Runtime**: Node.js 16+
- **Language**: TypeScript 4.9+
- **Web Framework**: Express.js
- **Database**: MongoDB (with Mongoose)
- **Cache**: Redis
- **Real-time**: Socket.IO
- **Queue**: Bull / BullMQ
- **Authentication**: JWT with cookie-session

### Additional Tools
- **File Upload**: Cloudinary SDK
- **Email**: SendGrid (production) + Nodemailer (development)
- **Testing**: Jest + ts-jest
- **Logging**: Bunyan
- **Validation**: Joi
- **Process Manager**: PM2
- **CI/CD**: CircleCI
- **Infrastructure**: Terraform (AWS)

---

## Features to Implement

Based on the Node.js reference implementation:

### Core Features
1. **Authentication**
   - Signup/Signin with JWT
   - Password reset flow with email
   - Session management

2. **User Management**
   - Profile CRUD operations
   - User search
   - Settings management

3. **Posts**
   - Create, read, update, delete posts
   - Image/video upload support
   - Post reactions (6 types: like, love, happy, wow, sad, angry)
   - Comments system

4. **Social Features**
   - Follow/unfollow users
   - Block/unblock users
   - Follower/following lists

5. **Real-time Chat**
   - Direct messaging
   - Message reactions
   - Read receipts
   - Typing indicators

6. **Notifications**
   - In-app notifications
   - Email notifications

7. **Media Upload**
   - Profile images
   - Background images
   - Post images/videos via Cloudinary

---

## Architecture Overview

The Node.js implementation will follow a feature-based modular architecture:

```
mysocial-backend/
├── src/
│   ├── features/         # Feature modules
│   │   ├── auth/        # Authentication
│   │   ├── user/        # User management
│   │   ├── post/        # Posts
│   │   ├── reactions/   # Reactions
│   │   ├── comments/    # Comments
│   │   ├── followers/   # Social graph
│   │   ├── chat/        # Messaging
│   │   ├── images/      # Images
│   │   └── notifications/ # Notifications
│   ├── shared/
│   │   ├── globals/     # Helpers, middleware
│   │   ├── services/    # DB, Cache, Queue services
│   │   ├── sockets/     # Socket.IO handlers
│   │   └── workers/     # Queue workers
│   ├── app.ts           # Entry point
│   ├── config.ts        # Configuration
│   ├── routes.ts        # Route aggregation
│   └── setupServer.ts   # Server setup
└── docs/                # Documentation
```

---

## Getting Started

### Prerequisites

- Node.js 16+ or higher
- MongoDB 5.0+
- Redis 7.0+
- Cloudinary account (free tier)
- Docker (optional)

### Setup Instructions

```bash
# Clone the repository
git clone -b node-master-onprogress git@github.com:kyimmQ/mysocial-backend.git
cd mysocial-backend

# Install dependencies
npm install

# Copy environment file
cp .env.development.example .env

# Update .env with your credentials
# DATABASE_URL, REDIS_HOST, CLOUD_NAME, etc.

# Run development server
npm run dev

# Server runs on http://localhost:5000
```

---

## Development Roadmap

### Phase 1: Foundation (Week 1-2)
- [ ] Project structure setup with TypeScript
- [ ] Configuration management
- [ ] MongoDB connection (Mongoose)
- [ ] Redis connection
- [ ] Express server setup
- [ ] Logging setup (Bunyan)
- [ ] Error handling middleware
- [ ] Security middleware (Helmet, CORS, HPP)

### Phase 2: Authentication (Week 3-4)
- [ ] User and Auth models (Mongoose schemas)
- [ ] Signup controller with validation (Joi)
- [ ] Login controller
- [ ] JWT middleware
- [ ] Password hashing (bcryptjs)
- [ ] Password reset flow with email
- [ ] Session management

### Phase 3: Core Features (Week 5-8)
- [ ] Post CRUD operations
- [ ] Post validation schemas
- [ ] Reaction system (6 types)
- [ ] Comment system
- [ ] User profile management
- [ ] Image upload (Cloudinary)
- [ ] Cache services (Redis)
- [ ] Queue setup (Bull)

### Phase 4: Social Features (Week 9-10)
- [ ] Follow/unfollow logic
- [ ] Block/unblock users
- [ ] Follower/following lists
- [ ] User search with regex
- [ ] Random user suggestions

### Phase 5: Real-time Features (Week 11-12)
- [ ] Socket.IO setup with Redis adapter
- [ ] Chat messaging
- [ ] Message reactions
- [ ] Read receipts
- [ ] Typing indicators
- [ ] Online user tracking

### Phase 6: Advanced Features (Week 13-14)
- [ ] Notification system
- [ ] Email notifications (SendGrid/Nodemailer)
- [ ] Background workers
- [ ] Complete caching strategy
- [ ] Bull Board UI
- [ ] API monitoring

### Phase 7: Testing & Deployment (Week 15-16)
- [ ] Unit tests (Jest)
- [ ] Controller tests
- [ ] Service tests
- [ ] Docker containerization
- [ ] CI/CD pipeline (CircleCI)
- [ ] AWS deployment (Terraform)

---

## Documentation

This repository includes comprehensive documentation:

- [Project Overview & PDR](./docs/project-overview-pdr.md)
- [System Architecture](./docs/system-architecture.md)
- [Codebase Summary](./docs/codebase-summary.md)
- [Code Standards](./docs/code-standards.md)
- [Caching Architecture](./docs/architecture/caching.md)
- [Queue Architecture](./docs/architecture/queues.md)
- [Project Roadmap](./docs/project-roadmap.md)

---

## Resources

### Node.js/TypeScript Learning Resources
- [Node.js Documentation](https://nodejs.org/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Express.js Guide](https://expressjs.com/en/guide/routing.html)

### Key Libraries
- [Express.js](https://expressjs.com/) - Web framework
- [Mongoose](https://mongoosejs.com/) - MongoDB ODM
- [Socket.IO](https://socket.io/docs/v4/) - Real-time communication
- [Bull](https://github.com/OptimalBits/bull) - Queue system
- [Jest](https://jestjs.io/) - Testing framework
- [Joi](https://joi.dev/) - Validation

---

## Reference Implementation

The complete Node.js/TypeScript implementation is available in the `node-main` branch for reference.

To see the complete implementation:
```bash
git checkout node-main
```

---

## Contributing

This is a personal learning project. Feel free to:
- Fork and implement your own version
- Submit issues for questions
- Share improvements and suggestions

---

## License

MIT License - Feel free to use for learning purposes.

---

**Happy Coding!** 🚀
