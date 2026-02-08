/**
 * ENVIRONMENT CONFIGURATION
 *
 * Validates and exports all environment variables.
 * Fails fast on startup if required vars are missing.
 *
 * Responsibilities:
 * - Load .env file (dotenv)
 * - Validate required variables (Joi or Zod)
 * - Export typed config object
 *
 * Example:
 *   export const config = {
 *     port: number,
 *     mongoUri: string,
 *     redisHost: string,
 *     jwtSecret: string,
 *     cloudinaryName: string,
 *     ...
 *   }
 */
