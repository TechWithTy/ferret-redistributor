package auth

// Permission literals used across the app.
const (
    PermOrgAdmin      = "org.admin"
    PermPostsRead     = "posts.read"
    PermPostsWrite    = "posts.write"
    PermScheduleClaim = "schedule.claim"
    PermScheduleManage= "schedule.manage"
    PermBillingRead   = "billing.read"
    PermBillingManage = "billing.manage"
    PermAnalyticsRead = "analytics.read"
)

// Role to permission mapping for quick checks (mirrors SQL seeds).
var RolePermissions = map[string][]string{
    "admin":  {PermOrgAdmin, PermPostsRead, PermPostsWrite, PermScheduleClaim, PermScheduleManage, PermBillingRead, PermBillingManage, PermAnalyticsRead},
    "editor": {PermPostsRead, PermPostsWrite, PermScheduleClaim, PermScheduleManage, PermAnalyticsRead},
    "viewer": {PermPostsRead, PermAnalyticsRead},
}

// HasPermission reports whether the given permission appears in the provided set.
func HasPermission(userPerms []string, want string) bool {
    for _, p := range userPerms {
        if p == want || p == PermOrgAdmin {
            return true
        }
    }
    return false
}

// ExpandRoles returns the union of permissions for role names.
func ExpandRoles(roles []string) []string {
    seen := map[string]struct{}{}
    out := make([]string, 0, 8)
    for _, r := range roles {
        for _, p := range RolePermissions[r] {
            if _, ok := seen[p]; ok { continue }
            seen[p] = struct{}{}
            out = append(out, p)
        }
    }
    return out
}

