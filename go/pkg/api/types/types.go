package types

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
}

type UserDTO struct {
    ID          string `json:"id"`
    Email       string `json:"email"`
    DisplayName string `json:"display_name"`
}

type SignupRequest struct {
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=8"`
    DisplayName string `json:"display_name" binding:"required"`
    OrgID       string `json:"org_id"`
}

type ForgotRequest struct {
    Email string `json:"email" binding:"required,email"`
}

type ResetRequest struct {
    Token       string `json:"token" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ProfileDTO struct {
    ID             string         `json:"id"`
    UserID         string         `json:"user_id"`
    OrgID          *string        `json:"org_id,omitempty"`
    DisplayName    *string        `json:"display_name,omitempty"`
    AvatarURL      *string        `json:"avatar_url,omitempty"`
    Bio            *string        `json:"bio,omitempty"`
    Timezone       *string        `json:"timezone,omitempty"`
    Locale         *string        `json:"locale,omitempty"`
    NotificationPrefs map[string]any `json:"notification_prefs,omitempty"`
    ContentPrefs   map[string]any `json:"content_prefs,omitempty"`
}

type ProfileUpdate struct {
    DisplayName *string        `json:"display_name"`
    AvatarURL   *string        `json:"avatar_url"`
    Bio         *string        `json:"bio"`
    Timezone    *string        `json:"timezone"`
    Locale      *string        `json:"locale"`
    NotificationPrefs map[string]any `json:"notification_prefs"`
    ContentPrefs map[string]any `json:"content_prefs"`
}

type ICPDTO struct {
    ID            string         `json:"id"`
    OrgID         string         `json:"org_id"`
    Name          string         `json:"name"`
    Timezone      *string        `json:"timezone,omitempty"`
    Languages     []string       `json:"languages,omitempty"`
    Region        *string        `json:"region,omitempty"`
    Industry      *string        `json:"industry,omitempty"`
    CompanySize   *string        `json:"company_size,omitempty"`
    Stage         *string        `json:"stage,omitempty"`
    Goals         map[string]any `json:"goals,omitempty"`
    Pains         map[string]any `json:"pains,omitempty"`
    BrandVoice    map[string]any `json:"brand_voice,omitempty"`
    Guidelines    *string        `json:"guidelines,omitempty"`
    Compliance    map[string]any `json:"compliance,omitempty"`
    Audience      map[string]any `json:"audience,omitempty"`
    ContentPillars map[string]any `json:"content_pillars,omitempty"`
}

type ICPUpdate struct {
    Name          *string        `json:"name"`
    Timezone      *string        `json:"timezone"`
    Languages     []string       `json:"languages"`
    Region        *string        `json:"region"`
    Industry      *string        `json:"industry"`
    CompanySize   *string        `json:"company_size"`
    Stage         *string        `json:"stage"`
    Goals         map[string]any `json:"goals"`
    Pains         map[string]any `json:"pains"`
    BrandVoice    map[string]any `json:"brand_voice"`
    Guidelines    *string        `json:"guidelines"`
    Compliance    map[string]any `json:"compliance"`
    Audience      map[string]any `json:"audience"`
    ContentPillars map[string]any `json:"content_pillars"`
}
