package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bitesinbyte/ferret/pkg/api/types"
	"github.com/gin-gonic/gin"
)

func getProfile(c *gin.Context) {
	uid := c.GetString(ctxUserID)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var p types.ProfileDTO
	var orgID sql.NullString
	var display, avatar, bio, tz, loc sql.NullString
	var notif, prefs []byte
	err := sqlDB.QueryRow(`SELECT id, user_id, org_id, display_name, avatar_url, bio, timezone, locale, notification_prefs, content_prefs FROM user_profiles WHERE user_id=$1`, uid).
		Scan(&p.ID, &p.UserID, &orgID, &display, &avatar, &bio, &tz, &loc, &notif, &prefs)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"profile": nil})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if orgID.Valid {
		p.OrgID = &orgID.String
	}
	if display.Valid {
		p.DisplayName = &display.String
	}
	if avatar.Valid {
		p.AvatarURL = &avatar.String
	}
	if bio.Valid {
		p.Bio = &bio.String
	}
	if tz.Valid {
		p.Timezone = &tz.String
	}
	if loc.Valid {
		p.Locale = &loc.String
	}
	_ = json.Unmarshal(notif, &p.NotificationPrefs)
	_ = json.Unmarshal(prefs, &p.ContentPrefs)
	c.JSON(http.StatusOK, p)
}

func updateProfile(c *gin.Context) {
	uid := c.GetString(ctxUserID)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req types.ProfileUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Upsert
	q := `INSERT INTO user_profiles(id, user_id, display_name, avatar_url, bio, timezone, locale, notification_prefs, content_prefs, created_at, updated_at)
          VALUES($1,$2,$3,$4,$5,$6,$7,COALESCE($8,'{}'::jsonb),COALESCE($9,'{}'::jsonb),NOW(),NOW())
          ON CONFLICT(user_id) DO UPDATE SET display_name=$3, avatar_url=$4, bio=$5, timezone=$6, locale=$7, notification_prefs=COALESCE($8,'{}'::jsonb), content_prefs=COALESCE($9,'{}'::jsonb), updated_at=NOW()`
	id := newID()
	b1, _ := json.Marshal(req.NotificationPrefs)
	b2, _ := json.Marshal(req.ContentPrefs)
	if _, err := sqlDB.Exec(q, id, uid, req.DisplayName, req.AvatarURL, req.Bio, req.Timezone, req.Locale, string(b1), string(b2)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
