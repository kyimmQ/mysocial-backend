// Package app wires all dependencies together (Dependency Injection).
//
// The Container holds all initialized dependencies and passes them
// to use-cases, handlers, and workers.
//
// Can use:
//   - Manual DI (recommended for small-medium projects)
//   - google/wire (compile-time DI code generation)
//
// Example:
//
//	type Container struct {
//	    Config      *config.Config
//	    DB          *mongo.Client
//	    Redis       *redis.Client
//	    UserHandler *handler.UserHandler
//	    AuthHandler *handler.AuthHandler
//	    // ...
//	}
//
//	func NewContainer(cfg *config.Config) (*Container, error) {
//	    // 1. Connect to MongoDB
//	    db, err := database.Connect(cfg.MongoURI)
//
//	    // 2. Connect to Redis
//	    rdb := cache.NewRedisClient(cfg.RedisHost)
//
//	    // 3. Build repositories (infrastructure implements domain interfaces)
//	    userRepo := database.NewMongoUserRepo(db)
//	    userCache := cache.NewRedisUserCache(rdb)
//
//	    // 4. Build use-cases (inject repositories)
//	    createUser := user.NewCreateUserUseCase(userRepo, userCache)
//	    getUser := user.NewGetUserUseCase(userRepo, userCache)
//
//	    // 5. Build handlers (inject use-cases)
//	    userHandler := handler.NewUserHandler(createUser, getUser)
//
//	    return &Container{...}, nil
//	}
package app
