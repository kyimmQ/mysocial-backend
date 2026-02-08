/**
 * CENTRAL ROUTE REGISTRY
 *
 * Aggregates all feature routes under a single base path.
 *
 * Example:
 *   app.use('/api/v1', authRoutes);
 *   app.use('/api/v1', postRoutes);
 *   app.use('/api/v1', chatRoutes);
 *   ...
 *
 * Benefits:
 * - Single place to see all endpoints
 * - Easy to add API versioning (/api/v1, /api/v2)
 * - Clean separation from Express app setup
 */
