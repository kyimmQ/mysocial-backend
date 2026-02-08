// Package repositories defines interfaces for data access.
// These are CONTRACTS — implementations live in infrastructure/database/.
//
// The domain layer owns these interfaces (Dependency Inversion Principle).
// Use-cases depend on these interfaces, NOT on concrete implementations.
//
// Example:
//
//	type UserRepository interface {
//	    Create(ctx context.Context, user *entities.User) error
//	    GetByID(ctx context.Context, id string) (*entities.User, error)
//	    GetByEmail(ctx context.Context, email string) (*entities.User, error)
//	    GetByUsername(ctx context.Context, username string) (*entities.User, error)
//	    Update(ctx context.Context, id string, user *entities.User) error
//	    Delete(ctx context.Context, id string) error
//	}
//
//	type AuthRepository interface {
//	    CreateAuth(ctx context.Context, auth *entities.Auth) error
//	    GetAuthByEmail(ctx context.Context, email string) (*entities.Auth, error)
//	    UpdatePassword(ctx context.Context, authID string, hash string) error
//	}
package repositories
