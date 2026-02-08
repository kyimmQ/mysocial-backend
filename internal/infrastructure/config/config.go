// Package config loads and validates environment configuration.
//
// Uses cleanenv, viper, or envconfig to:
//   1. Load .env file (development)
//   2. Read environment variables (production)
//   3. Validate required values (fail fast on startup)
//   4. Export a typed Config struct
//
// Example:
//
//	type Config struct {
//	    Port         int    `env:"PORT" env-default:"5000"`
//	    MongoURI     string `env:"MONGO_URI" env-required:"true"`
//	    RedisHost    string `env:"REDIS_HOST" env-required:"true"`
//	    JWTSecret    string `env:"JWT_SECRET" env-required:"true"`
//	    CloudName    string `env:"CLOUDINARY_NAME"`
//	    CloudKey     string `env:"CLOUDINARY_KEY"`
//	    CloudSecret  string `env:"CLOUDINARY_SECRET"`
//	    SendGridKey  string `env:"SENDGRID_API_KEY"`
//	    Environment  string `env:"NODE_ENV" env-default:"development"`
//	}
//
//	func Load() (*Config, error) {
//	    var cfg Config
//	    err := cleanenv.ReadEnv(&cfg)
//	    return &cfg, err
//	}
package config
