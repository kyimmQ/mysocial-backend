// Package user contains DTOs for the user use-case layer.
//
// DTOs (Data Transfer Objects) decouple use-case I/O from both
// domain entities and HTTP request/response shapes.
//
// Input DTOs: what the use-case receives
// Output DTOs: what the use-case returns
//
// Example:
//
//	type CreateUserInput struct {
//	    Username    string
//	    Email       string
//	    Password    string
//	    AvatarColor string
//	}
//
//	type CreateUserOutput struct {
//	    ID          string
//	    Username    string
//	    Email       string
//	    Token       string
//	    CreatedAt   time.Time
//	}
//
//	type GetUserOutput struct {
//	    ID             string
//	    Username       string
//	    Email          string
//	    PostsCount     int
//	    FollowersCount int
//	    FollowingCount int
//	}
package user
