package auth

import (
    "context"
    "database/sql"
)

// GetUserOrgPermissions returns the distinct permission names for a user within an org.
// It unions direct org-level role assignments and team-level assignments for teams in that org.
func GetUserOrgPermissions(ctx context.Context, db *sql.DB, orgID, userID string) ([]string, error) {
    const q = `
WITH perms AS (
    -- direct org user_roles
    SELECT rp.permission
    FROM user_roles ur
    JOIN role_permissions rp ON rp.role_id = ur.role_id
    WHERE ur.org_id = $1 AND ur.user_id = $2
    UNION
    -- team-scoped roles within same org
    SELECT rp.permission
    FROM team_roles tr
    JOIN teams t ON t.id = tr.team_id
    JOIN role_permissions rp ON rp.role_id = tr.role_id
    WHERE t.org_id = $1 AND tr.user_id = $2
)
SELECT DISTINCT permission FROM perms ORDER BY permission`

    rows, err := db.QueryContext(ctx, q, orgID, userID)
    if err != nil { return nil, err }
    defer rows.Close()
    out := make([]string, 0, 16)
    for rows.Next() {
        var p string
        if err := rows.Scan(&p); err != nil { return nil, err }
        out = append(out, p)
    }
    if err := rows.Err(); err != nil { return nil, err }
    return out, nil
}

// HasOrgPermission checks if user has a specific permission within an org.
func HasOrgPermission(ctx context.Context, db *sql.DB, orgID, userID, perm string) (bool, error) {
    const q = `
-- direct org roles
SELECT EXISTS (
    SELECT 1
    FROM user_roles ur
    JOIN role_permissions rp ON rp.role_id = ur.role_id
    WHERE ur.org_id = $1 AND ur.user_id = $2 AND (rp.permission = $3 OR rp.permission = 'org.admin')
) OR EXISTS (
    SELECT 1
    FROM team_roles tr
    JOIN teams t ON t.id = tr.team_id
    JOIN role_permissions rp ON rp.role_id = tr.role_id
    WHERE t.org_id = $1 AND tr.user_id = $2 AND (rp.permission = $3 OR rp.permission = 'org.admin')
)`
    var ok bool
    if err := db.QueryRowContext(ctx, q, orgID, userID, perm).Scan(&ok); err != nil {
        return false, err
    }
    return ok, nil
}

