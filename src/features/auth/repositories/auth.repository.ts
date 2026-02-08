/**
 * AUTH REPOSITORY
 *
 * Data access layer for auth feature.
 * Encapsulates all MongoDB/Redis operations for auth.
 *
 * Clean Architecture rule:
 *   Use-cases depend on repository INTERFACES, not implementations.
 *   This allows swapping MongoDB for PostgreSQL without touching business logic.
 *
 * Methods:
 *   - createAuthUser(data): Promise<IAuthDocument>
 *   - getAuthUserByEmail(email): Promise<IAuthDocument | null>
 *   - getAuthUserByUsername(username): Promise<IAuthDocument | null>
 *   - updatePassword(authId, hashedPassword): Promise<void>
 */
