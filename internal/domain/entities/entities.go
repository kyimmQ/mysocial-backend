// Package entities defines pure domain structs.
// These represent core business objects with NO framework tags,
// NO database annotations, NO JSON tags.
//
// Rules:
//   - Zero external imports (no mongo, no gorm, no gin)
//   - Only plain Go types
//   - Validation logic lives here if it's a business rule
//     (e.g., "email must not be empty" is domain, "email format" is pkg/validator)
//
// Examples:
//   - User: ID, Username, Email, Password, AvatarColor, CreatedAt
//   - Post: ID, UserID, Text, Privacy, Reactions, CommentsCount, CreatedAt
//   - Message: ID, SenderID, ReceiverID, Body, IsRead, CreatedAt
package entities
