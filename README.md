[![CircleCI](https://dl.circleci.com/status-badge/img/gh/uzochukwueddie/chatty-backend/tree/develop.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/uzochukwueddie/chatty-backend/tree/develop)
[![codecov](https://codecov.io/gh/uzochukwueddie/chatty-backend/branch/develop/graph/badge.svg?token=VR3XBTQMCV)](https://codecov.io/gh/uzochukwueddie/chatty-backend)

# MySocial Backend (Chatty App Backend)

|||||\n|:-:|:-:|:-:|:-:|\n|![First Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482458/github_readme_images/aws_bxdmec.png)|![Second Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482319/github_readme_images/Terraform_PrimaryLogo_Color_RGB_gcbknj.png)|![Third Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482279/github_readme_images/nodejs-logo_hqxxed.svg)|![Fourth Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482298/github_readme_images/ts-logo-512_jt9rmi.png)

|||||\n|:-:|:-:|:-:|:-:|\n|![First Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482275/github_readme_images/redis-icon_xzk6f2.png)|![Second Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482528/github_readme_images/Logo_RGB_Forest-Green_qjxd7x.png)|![Third Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482577/github_readme_images/pm2_owgicz.png)|![Fourth Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482745/github_readme_images/socketio_lcyu8y.jpg)

|||||\n|:-:|:-:|:-:|:-:|\n|![First Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482903/github_readme_images/Expressjs_sza4ue.png)|![Second Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662483106/github_readme_images/bull_y4erki.png)|![Third Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662482947/github_readme_images/sendgrid_d1v6dc.jpg)|![Fourth Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662483059/github_readme_images/nodemailer_rfpntx.png)

||\n|:-:|\n![First Image](https://res.cloudinary.com/dyamr9ym3/image/upload/v1662483242/github_readme_images/cloudinary_logo_blue_0720_2x_n8k46z.png)

## Overview

MySocial Backend is a production-ready, full-featured real-time social media backend application built with Node.js, TypeScript, Express, MongoDB, Redis, Socket.IO, and Bull queues. It powers a Twitter/Facebook-like social network platform with comprehensive features including posts, reactions, comments, real-time messaging, notifications, and social connections.

The system is designed for **high scalability**, **performance**, and **real-time user interactions** using asynchronous job processing, multi-layered caching, and WebSocket-based communication.

You can find the repo for the frontend built with React [here](https://github.com/uzochukwueddie/chatty).

---

## Quick Start

### Prerequisites

- Node.js 16.x or higher
- MongoDB 5.0+ (local or [MongoDB Atlas](https://www.mongodb.com/atlas/database))
- Redis 7.0+ (local or cloud)
- Cloudinary account ([sign up free](https://cloudinary.com/))
- Email service (Ethereal.email for dev, SendGrid for prod)

### Installation

```bash
# Clone the repository
git clone -b develop https://github.com/uzochukwueddie/chatty-backend.git
cd chatty-backend

# Install dependencies
npm install

# Create environment file
cp .env.development.example .env

# Update .env with your credentials
# DATABASE_URL, REDIS_HOST, CLOUD_NAME, etc.
```

### Running Locally

```bash
# Start MongoDB (if local)
mongod

# Start Redis (if local)
redis-server

# Start development server
npm run dev

# Server runs on http://localhost:5000
# Bull Board UI: http://localhost:5000/queues
```

**Important**: Inside `src/setupServer.ts`, comment the line `sameSite: 'none'` for local development. Uncomment before deploying to AWS.

### Running Tests

```bash
# Run unit tests with coverage
npm run test

# Run linting
npm run lint:check

# Run code formatting check
npm run prettier:check
```

### Seed Database (Optional)

```bash
# Generate fake data using Faker
npm run seeds:dev
```

---

## Architecture Overview

MySocial Backend follows a **layered architecture** with:

1. **Presentation Layer**: Express.js routes, controllers, middleware
2. **Business Logic Layer**: Feature modules, services, validation
3. **Data Access Layer**: MongoDB (Mongoose), Redis caching
4. **Queue & Worker Layer**: Bull queues for async processing
5. **Real-time Layer**: Socket.IO for live updates

**Key Patterns**:
- **Cache-first reads**: Check Redis before MongoDB (30-100x faster)
- **Write-through writes**: Save to cache immediately, queue DB writes asynchronously
- **Async processing**: Heavy operations (emails, media uploads) processed by workers
- **Real-time updates**: Socket.IO broadcasts changes to connected clients

**Request Flow**:
```
Client → Express → Controller → Cache (Redis) → Queue (Bull)
                                      ↓              ↓
                                   Response      Worker → MongoDB
                                      ↓
                                  Socket.IO → Real-time updates
```

For detailed architecture, see [System Architecture](./docs/system-architecture.md).

---

## Key Features

### 1. Authentication & User Management
- Signup and signin with JWT authentication
- Password reset via email (crypto tokens)
- Change password when logged in
- User profile management (basic info, social links, settings)
- User search and random suggestions
- Session management with encrypted cookies

### 2. Posts & Content
- Create, read, update, delete posts
- Support for text, images, videos, and GIFs (via Cloudinary)
- Post privacy settings (public, private, followers-only)
- Background colors for text posts
- Feelings/activity tags
- Paginated feeds with real-time updates

### 3. Social Interactions
- Six reaction types: like, love, happy, wow, sad, angry
- Comments on posts with notifications
- Follower/following system
- Block/unblock users
- Real-time reaction and comment updates via Socket.IO

### 4. Real-time Messaging
- Private chat with text, images, and GIFs
- Message reactions
- Read receipts (isRead flag)
- Soft delete (deleteForMe/deleteForEveryone)
- Typing indicators
- Conversation management

### 5. Notifications
- In-app notifications (follows, comments, reactions, messages)
- Email notifications (configurable per user)
- Mark as read/delete functionality
- Real-time notification delivery

### 6. Media Management
- Image uploads via Cloudinary
- Profile picture management
- Image gallery per user
- Background images for posts

---

## Technology Stack

### Backend
- **Node.js** v16+ - JavaScript runtime
- **TypeScript** v4.9+ - Type-safe JavaScript
- **Express.js** v4 - Web framework
- **Socket.IO** v4 - Real-time bidirectional communication

### Databases
- **MongoDB** v7 - Document database (Mongoose ODM)
- **Redis** v4+ - In-memory cache and session store

### Async Processing
- **Bull** v4 - Redis-based job queue
- **BullMQ** v3 - Modern queue system
- **Bull Board** v4 - Queue monitoring UI

### External Services
- **Cloudinary** - Image and video storage/transformation
- **SendGrid** - Transactional emails (production)
- **Nodemailer** - SMTP emails (development)

### DevOps & Infrastructure
- **PM2** - Production process manager (5 instances)
- **Terraform** - Infrastructure as code (AWS)
- **CircleCI** - CI/CD pipeline
- **AWS** - Cloud hosting (EC2, ALB, ElastiCache, Route53, S3)

### Testing & Code Quality
- **Jest** v29 - Testing framework
- **ESLint** v8 - Linting
- **Prettier** v2 - Code formatting
- **Faker** v7 - Test data generation

---

## Project Structure

```
mysocial-backend/
├── src/
│   ├── features/           # Feature modules (9 domains)
│   │   ├── auth/          # Authentication
│   │   ├── user/          # User profiles
│   │   ├── post/          # Posts
│   │   ├── reactions/     # Reactions
│   │   ├── comments/      # Comments
│   │   ├── followers/     # Social graph
│   │   ├── chat/          # Messaging
│   │   ├── images/        # Image management
│   │   └── notifications/ # Notifications
│   ├── shared/
│   │   ├── globals/       # Decorators, helpers, middleware
│   │   ├── services/      # Core services
│   │   │   ├── db/       # Database services (10)
│   │   │   ├── redis/    # Cache services (6)
│   │   │   ├── queues/   # Queue definitions (12)
│   │   │   └── emails/   # Email templates
│   │   ├── sockets/       # Socket.IO handlers (6)
│   │   └── workers/       # Queue workers (11)
│   ├── mocks/             # Test mocks
│   ├── app.ts             # Application entry point
│   ├── config.ts          # Configuration
│   ├── routes.ts          # Route aggregation
│   ├── seeds.ts           # Database seeding
│   ├── setupDatabase.ts   # MongoDB setup
│   └── setupServer.ts     # Express server setup
├── deployment/            # Terraform infrastructure
├── endpoints/             # API test files (.http)
├── scripts/               # Deployment scripts
└── docs/                  # Documentation
```

**Feature Module Pattern**:
```
feature/
├── controllers/   # Request handlers
├── models/        # Mongoose schemas
├── routes/        # Express routes
├── schemes/       # Joi validation schemas
└── interfaces/    # TypeScript types
```

---

## API Endpoints

The actual endpoints for the application can be found inside the `endpoints/` folder. The endpoint files all have a `.http` extension. To use these files to make API calls, install the [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension on VSCode.

**Available Endpoint Files**:
- `auth.http` - Authentication (signup, signin, password reset)
- `user.http` - User profiles, search
- `posts.http` - Post CRUD operations
- `reactions.http` - Reactions on posts
- `comments.http` - Comments on posts
- `follower.http` - Follow/unfollow, block/unblock
- `chat.http` - Real-time messaging
- `image.http` - Image uploads
- `notification.http` - Notifications
- `health.http` - Health checks

---

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[Project Overview & PDR](./docs/project-overview-pdr.md)** - Product requirements, features, technology rationale
- **[Codebase Summary](./docs/codebase-summary.md)** - Directory structure, module overview, code organization
- **[Code Standards](./docs/code-standards.md)** - TypeScript patterns, naming conventions, best practices
- **[System Architecture](./docs/system-architecture.md)** - High-level architecture, request flows
  - [Caching Architecture](./docs/architecture/caching.md) - Redis strategies, cache patterns
  - [Queue Architecture](./docs/architecture/queues.md) - Bull queues, async processing
- **[Project Roadmap](./docs/project-roadmap.md)** - Known issues, improvements, future plans

---

## Development Tools

### View Data

- **MongoDB**: Use [MongoDB Compass](https://www.mongodb.com/try/download/compass) to view database contents
- **Redis**: Use [Redis Commander](https://www.npmjs.com/package/redis-commander) to view cache contents
- **Queue Monitoring**: Access Bull Board UI at `http://localhost:5000/queues`

### Testing APIs

- **REST Client**: Use VSCode extension with `.http` files in `endpoints/` folder
- **Postman**: Import endpoints manually or use REST Client

---

## AWS Deployment

### Prerequisites

- AWS account with IAM permissions
- Domain name (for Route53 hosted zone)
- Terraform v1.0+ installed
- AWS CLI configured

### AWS Resources

The Terraform configuration creates:
- VPC with public/private subnets
- Internet Gateway and NAT Gateway
- Application Load Balancer (ALB)
- Auto Scaling Group (EC2 instances)
- ElastiCache (Redis cluster)
- Route53 DNS records
- SSL certificates (Certificate Manager)
- S3 buckets (env files, backups)
- CodeDeploy for automated deployments
- CloudWatch for monitoring

### Setup Steps

1. **Create Route53 Hosted Zone** (manual)
   - Create hosted zone on AWS Console
   - Copy NS records to your domain provider

2. **Create S3 Bucket for Environment Files**
   ```bash
   # Create bucket
   aws s3 mb s3://your-bucket-name

   # Create folder structure
   # your-bucket/backend/develop/

   # Zip and upload .env file
   zip env-file.zip .env.develop
   aws --region us-east-1 s3 cp env-file.zip s3://your-bucket/backend/develop/
   ```

3. **Update Terraform Variables**
   - Edit `deployment/variables.tf` with your values
   - Update `deployment/main.tf` with S3 bucket name
   - Add AWS keypair name to `ec2_launch_config.tf`

4. **Deploy Infrastructure**
   ```bash
   cd deployment
   terraform init
   terraform validate
   terraform fmt
   terraform plan
   terraform apply -auto-approve
   ```

5. **Destroy Resources** (when needed)
   ```bash
   terraform destroy
   ```

For detailed AWS setup, see the [AWS Setup](#aws-setup) section in the original README.

---

## CI/CD Pipeline with CircleCI

### Setup

1. Create account on [CircleCI](https://circleci.com/)
2. Connect your GitHub/Bitbucket repository
3. Add environment variables in CircleCI project settings:
   - `CODECOV_TOKEN` - From [CodeCov](https://about.codecov.io/)
   - `CODE_DEPLOY_UPDATE` - Set to `false` initially, `true` after first deploy
   - `SLACK_ACCESS_TOKEN` and `SLACK_DEFAULT_CHANNEL` - For build notifications

4. Update `.circleci/config.yml`:
   - Replace `<variable-prefix>` with your Terraform prefix value

### Pipeline Features

- Automated testing on every commit
- Code coverage reporting (CodeCov)
- Linting and formatting checks
- Automated deployment to AWS (CodeDeploy)
- Slack notifications for build status

---

## Monitoring

### Application Monitoring

- **Bull Board UI**: `http://localhost:5000/queues` - Queue monitoring
- **Swagger Stats**: Built-in API metrics (if enabled)
- **PM2 Monitoring**: `pm2 monit` - Process monitoring

### AWS Monitoring

- **CloudWatch**: Logs, metrics, alarms
- **ALB Metrics**: Request count, latency, errors
- **ElastiCache Metrics**: Redis performance

---

## Performance Characteristics

### Response Times

- **Cached reads**: 5-20ms (Redis)
- **Database reads**: 200-500ms (MongoDB)
- **Async writes**: 10-50ms (queue + cache)
- **Real-time events**: < 100ms (Socket.IO)

### Scalability

- **PM2 Clustering**: 5 instances by default
- **Queue Workers**: 60 concurrent jobs (12 queues × 5 workers)
- **Socket.IO**: Multi-instance support via Redis adapter
- **Auto Scaling**: Configured in AWS (min 2, max 10 instances)

### Cache Hit Rates

- **User profiles**: 80-90%
- **Posts feed**: 70-85%
- **Messages**: 60-75%

---

## Contributing

Contributions are welcome! Please see the [Project Roadmap](./docs/project-roadmap.md) for priority areas.

**Good First Issues**:
- Add cache TTL to Redis keys
- Add health check endpoints
- Write unit tests for services
- Implement password strength validation

**Development Process**:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following [Code Standards](./docs/code-standards.md)
4. Write tests for new features
5. Run tests and linting (`npm run test`, `npm run lint:check`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

---

## License

ISC

---

## Resources

- **Frontend Repository**: [Chatty Frontend](https://github.com/uzochukwueddie/chatty)
- **MongoDB**: [Installation Guide](https://www.mongodb.com/docs/manual/administration/install-community/)
- **Redis**: [Download](https://redis.io/download/)
- **Cloudinary**: [Sign Up](https://cloudinary.com/)
- **Ethereal Email**: [Testing Email](https://ethereal.email/)
- **SendGrid**: [Email Service](https://sendgrid.com/)

---

## Support

For questions or issues, please:
1. Check the [documentation](./docs/)
2. Review existing [issues](https://github.com/uzochukwueddie/chatty-backend/issues)
3. Create a new issue if needed

---

## Acknowledgments

Built with modern best practices for scalability, performance, and maintainability. Special thanks to all contributors and the open-source community.
