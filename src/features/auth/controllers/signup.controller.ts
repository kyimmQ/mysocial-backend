/**
 * AUTH CONTROLLERS
 *
 * THE BRIDGE between HTTP and business logic.
 *
 * Responsibilities:
 * 1. Extract data from req (body, params, query)
 * 2. Validate input (Joi/Zod schema)
 * 3. Call use-case
 * 4. Return HTTP response (status code + JSON)
 *
 * Controllers should be THIN — no business logic.
 * If a controller is getting complex, the logic belongs in a use-case.
 */
