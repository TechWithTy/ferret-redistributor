RBAC (Role-Based Access Control)

Tables
- roles: id, name, description, scope, is_system
- permissions: name, description
- role_permissions: role_id x permission
- user_roles: org_id x user_id x role_id

Seed Roles
- admin: full privileges (includes org.admin)
- editor: content/scheduling management
- viewer: read-only

Seed Permissions
- org.admin
- posts.read, posts.write
- schedule.claim, schedule.manage
- billing.read, billing.manage
- analytics.read

Integration
- App can query user roles for an org and expand to permissions.
- See pkg/auth/rbac.go for helpers and constants.

