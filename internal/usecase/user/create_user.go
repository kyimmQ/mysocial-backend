// Package user contains application use-cases for user operations.
//
// Each file = one action (Single Responsibility Principle):
//   - create_user.go  → signup logic
//   - get_user.go     → profile retrieval
//   - update_user.go  → profile updates
//   - delete_user.go  → account deletion
//
// Use-cases depend on domain/repositories interfaces (injected).
// They contain orchestration logic:
//   1. Validate business rules
//   2. Call repository methods
//   3. Return DTOs (not entities) to the interface layer
//
// Example:
//
//	type CreateUserUseCase struct {
//	    authRepo domain.AuthRepository
//	    userRepo domain.UserRepository
//	    cache    domain.UserCache
//	    queue    domain.JobQueue
//	}
//
//	func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
//	    // 1. Check if user exists
//	    // 2. Hash password
//	    // 3. Create auth + user records
//	    // 4. Cache user
//	    // 5. Queue background jobs
//	    // 6. Return output DTO
//	}
package user
