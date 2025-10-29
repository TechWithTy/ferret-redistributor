package recurpost

import "time"

// Common
type Pagination struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// Auth
type TokenRequest struct {
	GrantType    string `json:"grant_type"`              // password | refresh_token
	Username     string `json:"username,omitempty"`      // for password
	Password     string `json:"password,omitempty"`      // for password
	RefreshToken string `json:"refresh_token,omitempty"` // for refresh_token
}

// Posts
type CreatePostRequest struct {
	Text        string     `json:"text"`
	MediaIDs    []string   `json:"media_ids,omitempty"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	AccountIDs  []string   `json:"account_ids,omitempty"` // profiles to publish to
}

type UpdatePostRequest struct {
	Text        *string    `json:"text,omitempty"`
	MediaIDs    []string   `json:"media_ids,omitempty"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	Status      *string    `json:"status,omitempty"` // draft|scheduled|...
}

type ListPostsRequest struct {
	Pagination
	Status string `json:"status,omitempty"` // draft|scheduled|published|failed
}

// Media
type UploadMediaRequest struct {
	Filename string `json:"filename"`
	// Provide one of URL or Bytes (as base64 externally when wiring HTTP)
	URL   string `json:"url,omitempty"`
	Bytes []byte `json:"-"`
}

// User login
type UserLoginRequest struct {
	EmailID string `json:"emailid"`
	PassKey string `json:"pass_key"`
}

// Connect Social Account URLs
type ConnectSocialAccountURLsRequest struct {
	EmailID string `json:"emailid"`
	PassKey string `json:"pass_key"`
}

// Social account list
type SocialAccountListRequest struct {
	EmailID string `json:"emailid"`
}

// Library list
type LibraryListRequest struct {
	EmailID string `json:"emailid"`
	PassKey string `json:"pass_key"`
}

// History data
type HistoryDataRequest struct {
	EmailID           string `json:"emailid"`
	PassKey           string `json:"pass_key"`
	ID                string `json:"id,omitempty"`                   // social account id
	IsGetVideoUpdates string `json:"is_get_video_updates,omitempty"` // "true" or "false"
}

// Add content in library
type AddContentInLibraryRequest struct {
	EmailID            string   `json:"emailid"`
	PassKey            string   `json:"pass_key"`
	ID                 string   `json:"id"` // library id
	Message            string   `json:"message"`
	FBMessage          string   `json:"fb_message,omitempty"`
	TWMessage          string   `json:"tw_message,omitempty"`
	LNMessage          string   `json:"ln_message,omitempty"`
	INMessage          string   `json:"in_message,omitempty"`
	GMBMessage         string   `json:"gmb_message,omitempty"`
	PIMessage          string   `json:"pi_message,omitempty"`
	YTMessage          string   `json:"yt_message,omitempty"`
	TKMessage          string   `json:"tk_message,omitempty"`
	THMessage          string   `json:"th_message,omitempty"`
	BSMessage          string   `json:"bs_message,omitempty"`
	URL                string   `json:"url,omitempty"`
	ImageURL           []string `json:"image_url,omitempty"`
	VideoURL           string   `json:"video_url,omitempty"`
	FBPostType         string   `json:"fb_post_type,omitempty"` // feed|story|reel
	INPostType         string   `json:"in_post_type,omitempty"`
	INReelShareInFeed  string   `json:"in_reel_share_in_feed,omitempty"` // yes|no
	FirstComment       string   `json:"first_comment,omitempty"`
	FBFirstComment     string   `json:"fb_first_comment,omitempty"`
	LNFirstComment     string   `json:"ln_first_comment,omitempty"`
	INFirstComment     string   `json:"in_first_comment,omitempty"`
	LNDocument         string   `json:"ln_document,omitempty"`
	LNDocumentTitle    string   `json:"ln_document_title,omitempty"`
	PITitle            string   `json:"pi_title,omitempty"`
	GBPCTA             string   `json:"gbp_cta,omitempty"`
	GBPCTAURL          string   `json:"gbp_cta_url,omitempty"`
	GBPOfferTitle      string   `json:"gbp_offer_title,omitempty"`
	GBPOfferStartDate  string   `json:"gbp_offer_start_date,omitempty"`
	GBPOfferEndDate    string   `json:"gbp_offer_end_date,omitempty"`
	GBPOfferCouponCode string   `json:"gbp_offer_coupon_code,omitempty"`
	GBPOfferTerms      string   `json:"gbp_offer_terms,omitempty"`
	GBPRedeemOfferLink string   `json:"gbp_redeem_offer_link,omitempty"`
	YTTitle            string   `json:"yt_title,omitempty"`
	YTCategory         string   `json:"yt_category,omitempty"`
	YTPrivacyStatus    string   `json:"yt_privacy_status,omitempty"`
	YTUserTags         []string `json:"yt_user_tags,omitempty"`
}

// Post content on social account (ContentRequest + schedule_date_time)
type PostContentRequest struct {
	EmailID               string   `json:"emailid"`
	PassKey               string   `json:"pass_key"`
	ID                    string   `json:"id"` // social account id
	Message               string   `json:"message"`
	FBMessage             string   `json:"fb_message,omitempty"`
	TWMessage             string   `json:"tw_message,omitempty"`
	LNMessage             string   `json:"ln_message,omitempty"`
	INMessage             string   `json:"in_message,omitempty"`
	GMBMessage            string   `json:"gmb_message,omitempty"`
	PIMessage             string   `json:"pi_message,omitempty"`
	YTMessage             string   `json:"yt_message,omitempty"`
	TKMessage             string   `json:"tk_message,omitempty"`
	THMessage             string   `json:"th_message,omitempty"`
	BSMessage             string   `json:"bs_message,omitempty"`
	URL                   string   `json:"url,omitempty"`
	ImageURL              []string `json:"image_url,omitempty"`
	VideoURL              string   `json:"video_url,omitempty"`
	FBPostType            string   `json:"fb_post_type,omitempty"` // feed|story|reel
	INPostType            string   `json:"in_post_type,omitempty"`
	INReelShareInFeed     string   `json:"in_reel_share_in_feed,omitempty"` // yes|no
	FirstComment          string   `json:"first_comment,omitempty"`
	FBFirstComment        string   `json:"fb_first_comment,omitempty"`
	LNFirstComment        string   `json:"ln_first_comment,omitempty"`
	INFirstComment        string   `json:"in_first_comment,omitempty"`
	LNDocument            string   `json:"ln_document,omitempty"`
	LNDocumentTitle       string   `json:"ln_document_title,omitempty"`
	PITitle               string   `json:"pi_title,omitempty"`
	GBPCTA                string   `json:"gbp_cta,omitempty"`
	GBPCTAURL             string   `json:"gbp_cta_url,omitempty"`
	GBPOfferTitle         string   `json:"gbp_offer_title,omitempty"`
	GBPOfferStartDate     string   `json:"gbp_offer_start_date,omitempty"`
	GBPOfferEndDate       string   `json:"gbp_offer_end_date,omitempty"`
	GBPOfferCouponCode    string   `json:"gbp_offer_coupon_code,omitempty"`
	GBPOfferTerms         string   `json:"gbp_offer_terms,omitempty"`
	GBPRedeemOfferLink    string   `json:"gbp_redeem_offer_link,omitempty"`
	YTTitle               string   `json:"yt_title,omitempty"`
	YTCategory            string   `json:"yt_category,omitempty"`
	YTPrivacyStatus       string   `json:"yt_privacy_status,omitempty"`
	YTUserTags            []string `json:"yt_user_tags,omitempty"`
	YTThumb               string   `json:"yt_thumb,omitempty"`
	YTMadeForKids         string   `json:"yt_video_made_for_kids,omitempty""`
	TKPrivacyStatus       string   `json:"tk_privacy_status,omitempty"`
	TKAllowComments       string   `json:"tk_allow_comments,omitempty"`
	TKAllowDuet           string   `json:"tk_allow_duet,omitempty"`
	TKAllowStitches       string   `json:"tk_allow_stitches,omitempty"`
	IsTopOfQueue          *int     `json:"is_top_of_queue,omitempty"`
	HostImagesOnRecurPost string   `json:"host_images_on_recurpost,omitempty"`
	ContentLiveDate       string   `json:"content_livedate,omitempty"`
	ContentExpireDate     string   `json:"content_expiredate,omitempty"`
	ScheduleDateTime      string   `json:"schedule_date_time,omitempty"`
}

// Generate content with AI
type GenerateContentWithAIRequest struct {
	EmailID      string `json:"emailid"`
	PassKey      string `json:"pass_key"`
	PromptText   string `json:"prompt_text"`
	AIID         string `json:"ai_id,omitempty"`
	ChatProgress string `json:"chat_progress,omitempty"`
	ChatHistory  string `json:"chat_history,omitempty"`
}

// Generate image with AI
type GenerateImageWithAIRequest struct {
	EmailID    string `json:"emailid"`
	PassKey    string `json:"pass_key"`
	PromptText string `json:"prompt_text"`
}
