/**
 * MAILER ADAPTER
 *
 * Abstracts email sending behind a clean interface.
 *
 * Implementations:
 * - Development: Nodemailer with Ethereal (fake SMTP)
 * - Production: SendGrid API
 *
 * Interface:
 *   sendEmail(to: string, subject: string, html: string): Promise<void>
 *
 * Templates live alongside features or in shared/globals/.
 */
