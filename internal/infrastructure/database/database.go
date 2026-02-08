// Package database handles database connections and repository implementations.
//
// Connection lifecycle:
//   - Connect(ctx, uri) → *mongo.Client
//   - Disconnect(ctx) → close connections
//   - Health check
//
// Repository implementations live alongside:
//   - user_repo_mongo.go  → implements domain.UserRepository using MongoDB
//   - auth_repo_mongo.go  → implements domain.AuthRepository using MongoDB
//
// These files contain Mongo-specific structs with bson tags,
// mapping between domain entities and database documents.
//
// Example:
//
//	type mongoUserRepo struct {
//	    collection *mongo.Collection
//	}
//
//	func (r *mongoUserRepo) GetByID(ctx context.Context, id string) (*entities.User, error) {
//	    var doc userDocument  // internal struct with bson tags
//	    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
//	    return doc.toEntity(), err
//	}
package database
