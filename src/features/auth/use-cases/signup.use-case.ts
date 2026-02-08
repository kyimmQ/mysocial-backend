/**
 * AUTH USE-CASES
 *
 * PURE BUSINESS LOGIC — no Express, no MongoDB, no Redis here.
 *
 * Each use-case is a single action:
 *   - SignUp: validate uniqueness, hash password, generate JWT, create user
 *   - SignIn: find user, compare password, generate JWT
 *   - ForgotPassword: generate reset token, send email
 *   - ResetPassword: verify token, update password
 *
 * Dependencies are INJECTED (repository interfaces, not concrete classes).
 * This makes use-cases testable with simple mocks.
 */
