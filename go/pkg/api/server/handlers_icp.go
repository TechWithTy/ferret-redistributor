package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bitesinbyte/ferret/pkg/api/types"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func getICP(c *gin.Context) {
	uid := c.GetString(ctxUserID)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Derive org_id from profile
	var orgID sql.NullString
	_ = sqlDB.QueryRow(`SELECT org_id FROM user_profiles WHERE user_id=$1`, uid).Scan(&orgID)
	if !orgID.Valid {
		c.JSON(http.StatusOK, gin.H{"icp": nil})
		return
	}
	// Return first profile by name
	var dto types.ICPDTO
	var languages []string
	var goals, pains, brand, compliance, audience, pillars []byte
	err := sqlDB.QueryRow(`SELECT id, org_id, name, timezone, languages, region, industry, company_size, stage, goals, pains, brand_voice, guidelines, compliance, audience, content_pillars FROM icp_profiles WHERE org_id=$1 ORDER BY name ASC LIMIT 1`, orgID.String).
		Scan(&dto.ID, &dto.OrgID, &dto.Name, &dto.Timezone, pq.Array(&languages), &dto.Region, &dto.Industry, &dto.CompanySize, &dto.Stage, &goals, &pains, &brand, &dto.Guidelines, &compliance, &audience, &pillars)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"icp": nil})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dto.Languages = languages
	_ = json.Unmarshal(goals, &dto.Goals)
	_ = json.Unmarshal(pains, &dto.Pains)
	_ = json.Unmarshal(brand, &dto.BrandVoice)
	_ = json.Unmarshal(compliance, &dto.Compliance)
	_ = json.Unmarshal(audience, &dto.Audience)
	_ = json.Unmarshal(pillars, &dto.ContentPillars)
	c.JSON(http.StatusOK, dto)
}

func updateICP(c *gin.Context) {
	uid := c.GetString(ctxUserID)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var orgID string
	if err := sqlDB.QueryRow(`SELECT COALESCE(org_id,'') FROM user_profiles WHERE user_id=$1`, uid).Scan(&orgID); err != nil || orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no org associated with profile"})
		return
	}
	var req types.ICPUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Upsert by (org_id, name='Default')
	name := "Default"
	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	}
	gb, _ := json.Marshal(req.Goals)
	pb, _ := json.Marshal(req.Pains)
	bb, _ := json.Marshal(req.BrandVoice)
	cb, _ := json.Marshal(req.Compliance)
	ab, _ := json.Marshal(req.Audience)
	pb2, _ := json.Marshal(req.ContentPillars)
	q := `INSERT INTO icp_profiles(id, org_id, name, timezone, languages, region, industry, company_size, stage, goals, pains, brand_voice, guidelines, compliance, audience, content_pillars, created_at, updated_at)
          VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,COALESCE($10,'{}'::jsonb),COALESCE($11,'{}'::jsonb),COALESCE($12,'{}'::jsonb),$13,COALESCE($14,'{}'::jsonb),COALESCE($15,'{}'::jsonb),COALESCE($16,'{}'::jsonb),NOW(),NOW())
          ON CONFLICT(org_id, name) DO UPDATE SET timezone=$4, languages=$5, region=$6, industry=$7, company_size=$8, stage=$9, goals=COALESCE($10,'{}'::jsonb), pains=COALESCE($11,'{}'::jsonb), brand_voice=COALESCE($12,'{}'::jsonb), guidelines=$13, compliance=COALESCE($14,'{}'::jsonb), audience=COALESCE($15,'{}'::jsonb), content_pillars=COALESCE($16,'{}'::jsonb), updated_at=NOW()`
	id := newID()
	if _, err := sqlDB.Exec(q, id, orgID, name, req.Timezone, pq.Array(req.Languages), req.Region, req.Industry, req.CompanySize, req.Stage, string(gb), string(pb), string(bb), req.Guidelines, string(cb), string(ab), string(pb2)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
