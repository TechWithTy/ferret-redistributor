# Database

Migrations live under `_data/_db/**` and are applied in lexical order by the migration runner.

Core tables
- Organizations, users, teams
- Social accounts, content items, campaigns
- scheduled_posts (Go calendar consumer)

AI & Optimization
- ai_generations, ai_variants, experiments, experiment_arms, post_outcomes
- trend_metrics (analytics), app_metrics (telemetry fallback)

Billing
- pricing_plans, org_subscriptions
- credit_wallets, credit_transactions
- success_metrics, success_tiers, org_success_enrollments

Admin & Support
- product_events, admin_kpis
- admin_profiles, support_queues, support_tickets, support_messages

Auth
- auth_identities (email/phone/linkedin/meta), auth_sessions
- auth_email_verifications, auth_password_resets, auth_phone_codes, auth_logins

User Personalization
- user_profiles
- icp_profiles + related tables (platform preferences, competitors, constraints)

