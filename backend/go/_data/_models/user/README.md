User Authentication (Email, Phone, LinkedIn, Meta)

Tables
- auth_identities: provider= email | phone | linkedin | meta; identifier (email/phone/provider id), secret_hash for email/password, oauth_data JSONB, verified_at, is_primary.
- auth_email_verifications: one-time tokens to verify email addresses.
- auth_password_resets: one-time tokens to reset passwords.
- auth_phone_codes: short-lived hashed OTP codes to verify phone numbers.
- auth_sessions: opaque session tokens (store hash only), expirable, with audit fields.
- auth_logins: audit trail of login attempts.

Flows (MVP)
- Sign up (email): create user + auth_identity(email), send email verification token; on verify, set verified_at.
- Password login: compare bcrypt hash in secret_hash (store only bcrypt hash), on success create auth_session.
- Forgot password: create reset token in auth_password_resets; on redeem and verify, rotate secret_hash.
- Phone verification: create auth_phone_codes, send via SMS, verify and set verified_at on phone identity.
- OAuth (LinkedIn/Meta): create identity with provider id in identifier, store tokens/claims in oauth_data. No secret_hash.

Indexes & Constraints
- Unique (provider, identifier) ensures no duplicates across users.
- Sessions indexed by user_id and expires_at for cleanup.
- Store tokens and OTPs hashed; never store raw tokens.

