package mock

import (
	"math/rand"
	"time"

	"github.com/bitesinbyte/ferret/pkg/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

// MockDataSuite contains all the mock data for testing
type MockDataSuite struct {
	Organizations       []models.Organization
	Users               []models.User
	Teams               []models.Team
	TeamMembers         []models.TeamMember
	SocialAccounts      []models.SocialAccount
	ContentItems        []models.ContentItem
	Campaigns           []models.Campaign
	ScheduledPosts      []models.ScheduledPost
	AIGenerations       []models.AIGeneration
	AIVariants          []models.AIVariant
	Experiments         []models.Experiment
	ExperimentArms      []models.ExperimentArm
	PostOutcomes        []models.PostOutcome
	TrendMetrics        []models.TrendMetric
	AppMetrics          []models.AppMetric
	MarketplacePosts    []models.MarketplacePost
	MarketplaceTxns     []models.MarketplaceTransaction
	ContentSources      []models.ContentSource
	ContentSourceItems  []models.ContentSourceItem
	YouTubePlaylists    []models.YouTubePlaylist
	ContentSyncs        []models.ContentSourceSync
	OAuthProviders      []models.OAuthProvider
	OAuthAccounts       []models.OAuthAccount
	PhoneVerifications  []models.PhoneVerification
	SocialAccountLinks  []models.SocialAccountLink
}

// NewMockDataSuite creates a new instance of MockDataSuite with populated mock data
func NewMockDataSuite() *MockDataSuite {
	suite := &MockDataSuite{}
	suite.generateMockData()
	return suite
}

func (m *MockDataSuite) generateMockData() {
	// Generate organizations
	orgs := GenerateMockOrganizations(3)
	m.Organizations = orgs
	
	// Generate OAuth providers
	m.OAuthProviders = GenerateMockOAuthProviders()
	
	// Generate users for each organization
	for _, org := range orgs {
		users := GenerateMockUsers(5, org.ID)
		m.Users = append(m.Users, users...)

		// Generate teams for each organization
		teams := GenerateMockTeams(2, org.ID)
		m.Teams = append(m.Teams, teams...)

		// Generate team members
		for _, team := range teams {
			userIDs := make([]string, 0, 3)
			for i := 0; i < 3 && i < len(users); i++ {
				userIDs = append(userIDs, users[i].ID)
			}
			members := GenerateMockTeamMembers(team.ID, userIDs)
			m.TeamMembers = append(m.TeamMembers, members...)

			// Generate social accounts for the organization
			socialAccounts := GenerateMockSocialAccounts(3, org.ID, nil)
			m.SocialAccounts = append(m.SocialAccounts, socialAccounts...)
		
			// Generate OAuth accounts for users
			oauthAccounts := GenerateMockOAuthAccounts(users, m.OAuthProviders)
			m.OAuthAccounts = append(m.OAuthAccounts, oauthAccounts...)
		
			// Generate social account links
			socialAccountLinks := GenerateMockSocialAccountLinks(users)
			m.SocialAccountLinks = append(m.SocialAccountLinks, socialAccountLinks...)
		}

		// Generate content items
		contentItems := GenerateMockContentItems(5, org.ID, users[0].ID)
		m.ContentItems = append(m.ContentItems, contentItems...)

		// Generate campaigns
		campaigns := GenerateMockCampaigns(2, org.ID)
		m.Campaigns = append(m.Campaigns, campaigns...)

		// Generate scheduled posts
		for _, account := range m.SocialAccounts {
			if account.OrgID == org.ID {
				for i := 0; i < 3; i++ {
					var campaignID *string
					if i%2 == 0 && len(campaigns) > 0 {
						campaignID = &campaigns[0].ID
					}

					var contentID *string
					if len(contentItems) > 0 {
						contentID = &contentItems[i%len(contentItems)].ID
					}

					posts := GenerateMockScheduledPosts(2, org.ID, campaignID, contentID, &account.ID)
					m.ScheduledPosts = append(m.ScheduledPosts, posts...)
				}
			}
		}

		// Generate AI generations and variants
		for _, user := range users {
			for i := 0; i < 2; i++ {
				gen := GenerateMockAIGeneration(org.ID, &user.ID, nil)
				m.AIGenerations = append(m.AIGenerations, gen)

				// Generate variants for each generation
				variants := GenerateMockAIVariants(gen.ID, 3)
				m.AIVariants = append(m.AIVariants, variants...)
			}
		}

		// Generate experiments
		experiments := GenerateMockExperiments(2, org.ID)
		m.Experiments = append(m.Experiments, experiments...)

		// Generate experiment arms
		for _, exp := range experiments {
			arms := GenerateMockExperimentArms(exp.ID, 3)
			m.ExperimentArms = append(m.ExperimentArms, arms...)
		}

		// Generate post outcomes for scheduled posts
		for _, post := range m.ScheduledPosts {
			if post.OrgID == org.ID && post.Status == "published" && post.PublishedAt != nil {
				outcome := GenerateMockPostOutcome(post.ID, post.Platform)
				m.PostOutcomes = append(m.PostOutcomes, outcome)
			}
		}

		// Generate trend metrics
		m.TrendMetrics = append(m.TrendMetrics, GenerateMockTrendMetrics(org.ID, 10)...)

		// Generate app metrics
		m.AppMetrics = GenerateMockAppMetrics(20)
	
		// Generate phone verifications
		m.PhoneVerifications = GenerateMockPhoneVerifications(10)

		// Generate marketplace posts and transactions
		for i := 0; i < 3; i++ {
			if len(users) > 0 && len(contentItems) > i {
				post := GenerateMockMarketplacePost(org.ID, users[0].ID, &contentItems[i].ID)
				m.MarketplacePosts = append(m.MarketplacePosts, post)

				if i < len(users)-1 {
					tx := GenerateMockMarketplaceTransaction(post.ID, users[1].ID)
					m.MarketplaceTxns = append(m.MarketplaceTxns, tx)
				}
			}
		}
	}
}

// GenerateMockContentItems generates mock content items
func GenerateMockContentItems(count int, orgID, userID string) []models.ContentItem {
	var items []models.ContentItem

	for i := 0; i < count; i++ {
		item := models.ContentItem{
			ID:           uuid.New().String(),
			OrgID:        orgID,
			Title:        faker.Sentence(),
			Body:         faker.Paragraph(),
			CanonicalURL: faker.URL(),
			MediaURL:     faker.URL() + "/image.jpg",
			Metadata: map[string]interface{}{
				"author":    faker.Name(),
				"wordCount": rand.Intn(1901) + 100, // Random between 100-2000
				"tags":      []string{faker.Word(), faker.Word(), faker.Word()},
			},
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
		items = append(items, item)
	}

	return items
}

// GenerateMockCampaigns generates mock campaigns
func GenerateMockCampaigns(count int, orgID string) []models.Campaign {
	var campaigns []models.Campaign

	for i := 0; i < count; i++ {
		campaign := models.Campaign{
			ID:          uuid.New().String(),
			OrgID:       orgID,
			Name:        "Campaign " + faker.Word() + " " + faker.Word(),
			Description: faker.Sentence(),
			CreatedAt:   time.Now().Add(-time.Duration(i) * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-time.Duration(i) * 24 * time.Hour),
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns
}

// GenerateMockScheduledPosts generates mock scheduled posts
func GenerateMockScheduledPosts(count int, orgID string, campaignID, contentID, socialAccountID *string) []models.ScheduledPost {
	var posts []models.ScheduledPost
	statuses := []string{"draft", "scheduled", "published", "failed"}

	for i := 0; i < count; i++ {
		status := statuses[i%len(statuses)]
		var publishedAt *time.Time
		if status == "published" {
			t := time.Now().Add(-time.Duration(i) * time.Hour)
			publishedAt = &t
		}

		post := models.ScheduledPost{
			ID:              uuid.New().String(),
			OrgID:           orgID,
			CampaignID:      campaignID,
			ContentID:       contentID,
			SocialAccountID: socialAccountID,
			Platform:        []string{"instagram", "twitter", "linkedin"}[i%3],
			Caption:         stringPtr(faker.Sentence() + " " + faker.Word()),
			Hashtags:        stringPtr("#" + faker.Word() + " #" + faker.Word()),
			ScheduledAt:     time.Now().Add(time.Duration(i+1) * time.Hour),
			Status:          status,
			ExternalID:      stringPtr(uuid.New().String()),
			PublishedAt:     publishedAt,
			Metadata: map[string]interface{}{
				"retry_count": i % 3,
			},
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
		posts = append(posts, post)
	}

	return posts
}

// GenerateMockAIGeneration generates a mock AI generation
func GenerateMockAIGeneration(orgID string, userID, contentItemID *string) models.AIGeneration {
	return models.AIGeneration{
		ID:            uuid.New().String(),
		OrgID:         orgID,
		UserID:        userID,
		Model:         "gpt-4",
		Prompt:        "Generate a blog post about " + faker.Word() + " and " + faker.Word(),
		Parameters:    map[string]interface{}{"temperature": 0.7, "max_tokens": 1000},
		OutputText:    faker.Paragraph(),
		OutputJSON:    map[string]interface{}{"title": faker.Sentence(), "body": faker.Paragraph()},
		ContentItemID: contentItemID,
		Status:        "completed",
		ErrorMessage:  "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// GenerateMockAIVariants generates mock AI variants
func GenerateMockAIVariants(generationID string, count int) []models.AIVariant {
	var variants []models.AIVariant

	for i := 0; i < count; i++ {
		variant := models.AIVariant{
			ID:           uuid.New().String(),
			GenerationID: generationID,
			VariantIndex: i,
			Payload: map[string]interface{}{
				"title":   faker.Sentence(),
				"content": faker.Paragraph(),
				"tone":    []string{"professional", "casual", "persuasive", "informative"}[i%4],
			},
			CreatedAt: time.Now().Add(time.Duration(i) * time.Minute),
		}
		variants = append(variants, variant)
	}

	return variants
}

// GenerateMockExperiments generates mock experiments
func GenerateMockExperiments(count int, orgID string) []models.Experiment {
	var experiments []models.Experiment

	for i := 0; i < count; i++ {
		endedAt := time.Now().Add(time.Duration(i+1) * 24 * time.Hour)
		experiment := models.Experiment{
			ID:          uuid.New().String(),
			OrgID:       orgID,
			Name:        "Experiment " + faker.Word(),
			Hypothesis:  "If we " + faker.Word() + " then we will see an increase in " + faker.Word(),
			StartedAt:   time.Now().Add(-time.Duration(i) * 24 * time.Hour),
			EndedAt:     &endedAt,
			CreatedAt:   time.Now().Add(-time.Duration(i+1) * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-time.Duration(i) * 24 * time.Hour),
		}
		experiments = append(experiments, experiment)
	}

	return experiments
}

// GenerateMockExperimentArms generates mock experiment arms
func GenerateMockExperimentArms(experimentID string, count int) []models.ExperimentArm {
	var arms []models.ExperimentArm

	for i := 0; i < count; i++ {
		arm := models.ExperimentArm{
			ID:           uuid.New().String(),
			ExperimentID: experimentID,
			VariantID:    stringPtr(uuid.New().String()),
			Weight:       float64(1) / float64(count),
			CreatedAt:    time.Now(),
		}
		arms = append(arms, arm)
	}

	return arms
}

// GenerateMockPostOutcome generates a mock post outcome
func GenerateMockPostOutcome(scheduledPostID, platform string) models.PostOutcome {
	return models.PostOutcome{
		ID:             uuid.New().String(),
		ScheduledPostID: scheduledPostID,
		Platform:       platform,
		ExternalID:     uuid.New().String(),
		Impressions:    int64(rand.Intn(9901) + 100),     // 100-10000
		Reach:          int64(rand.Intn(7921) + 80),       // 80-8000
		Likes:          int64(rand.Intn(991) + 10),        // 10-1000
		Comments:       int64(rand.Intn(101)),             // 0-100
		Shares:         int64(rand.Intn(501)),             // 0-500
		Clicks:         int64(rand.Intn(496) + 5),         // 5-500
		Saves:          int64(rand.Intn(201)),             // 0-200
		Conversions:    int64(rand.Intn(51)),              // 0-50
		CollectedAt:    time.Now(),
		Metadata: map[string]interface{}{
			"engagement_rate": 0.5 + rand.Float64()*9.5, // 0.5-10.0
		},
	}
}

// GenerateMockTrendMetrics generates mock trend metrics
func GenerateMockTrendMetrics(orgID string, count int) []models.TrendMetric {
	var metrics []models.TrendMetric
	sources := []string{"instagram", "twitter", "linkedin", "website", "newsletter"}
	dimensions := []string{"impressions", "engagement", "clicks", "conversions"}

	for i := 0; i < count; i++ {
		source := sources[i%len(sources)]
		dimension := dimensions[i%len(dimensions)]
		bucketStart := time.Now().Add(-time.Duration(24*(count-i)) * time.Hour)

		metric := models.TrendMetric{
			OrgID:       orgID,
			Source:      source,
			Dimension:   dimension,
			Metric:      "count",
			BucketStart: bucketStart,
			BucketEnd:   bucketStart.Add(24 * time.Hour),
			Value:       float64(rand.Intn(9901) + 100), // 100-10000
			Meta: map[string]interface{}{
				"platform": source,
			},
			CreatedAt: time.Now(),
		}
		metrics = append(metrics, metric)
	}

	return metrics
}

// GenerateMockAppMetrics generates mock app metrics
func GenerateMockAppMetrics(count int) []models.AppMetric {
	var metrics []models.AppMetric
	names := []string{"api.latency", "db.queries", "cache.hits", "cache.misses", "users.active"}

	for i := 0; i < count; i++ {
		metric := models.AppMetric{
			Name:  names[i%len(names)],
			Value: rand.Float64() * 1000, // 0-1000
			Attributes: map[string]interface{}{
				"status":     "success",
				"endpoint":   "/api/" + faker.Word(),
				"http_method": []string{"GET", "POST", "PUT", "DELETE"}[i%4],
			},
			RecordedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
		metrics = append(metrics, metric)
	}

	return metrics
}

// GenerateMockMarketplacePost generates a mock marketplace post
func GenerateMockMarketplacePost(orgID, sellerUserID string, contentItemID *string) models.MarketplacePost {
	return models.MarketplacePost{
		ID:            uuid.New().String(),
		OrgID:         orgID,
		SellerUserID:  sellerUserID,
		ContentItemID: contentItemID,
		Title:         "Content: " + faker.Sentence(),
		Description:   faker.Paragraph(),
		PriceCents:    rand.Intn(4501) + 500, // 500-5000
		Currency:      "USD",
		Status:        []string{"draft", "published", "sold"}[rand.Intn(3)],
		Metadata: map[string]interface{}{
			"tags":      []string{faker.Word(), faker.Word()},
			"category":  []string{"article", "image", "video", "template"}[rand.Intn(4)],
			"wordCount": rand.Intn(1901) + 100, // 100-2000
		},
		CreatedAt: time.Now().Add(-time.Duration(rand.Intn(30)+1) * 24 * time.Hour), // 1-30 days ago
		UpdatedAt: time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),    // 0-29 days ago
	}
}

// GenerateMockMarketplaceTransaction generates a mock marketplace transaction
func GenerateMockMarketplaceTransaction(postID, buyerUserID string) models.MarketplaceTransaction {
	return models.MarketplaceTransaction{
		ID:          uuid.New().String(),
		PostID:      postID,
		BuyerUserID: buyerUserID,
		AmountCents: rand.Intn(4501) + 500, // 500-5000
		Currency:    "USD",
		Status:      "completed",
		CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(15)) * 24 * time.Hour), // 0-14 days ago
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
