// Package cache handles Redis connections and cache operations.
//
// Provides:
//   - Redis client initialization (go-redis)
//   - Cache interface implementations
//   - Data structure operations (HASH, ZSET, LIST)
//
// Cache implementations satisfy domain-level cache interfaces,
// keeping use-cases decoupled from Redis specifics.
//
// Example:
//
//	type redisUserCache struct {
//	    client *redis.Client
//	}
//
//	func (c *redisUserCache) SaveUser(ctx context.Context, user *entities.User) error {
//	    data, _ := json.Marshal(user)
//	    return c.client.HSet(ctx, "users", user.ID, data).Err()
//	}
//
//	func (c *redisUserCache) GetUser(ctx context.Context, id string) (*entities.User, error) {
//	    data, err := c.client.HGet(ctx, "users", id).Result()
//	    // unmarshal and return
//	}
package cache
