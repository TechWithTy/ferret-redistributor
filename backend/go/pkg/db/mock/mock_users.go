package mock

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bitesinbyte/ferret/pkg/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

// GenerateMockUsers generates a slice of mock User models
func GenerateMockUsers(count int, orgID string) []models.User {
	var users []models.User

	for i := 0; i < count; i++ {
		user := models.User{
			ID:          uuid.New().String(),
			OrgID:       orgID,
			Email:       faker.Email(),
			DisplayName: faker.Name(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		users = append(users, user)
	}

	return users
}

// GenerateMockOrganizations generates a slice of mock Organization models
func GenerateMockOrganizations(count int) []models.Organization {
	var orgs []models.Organization

	for i := 0; i < count; i++ {
		org := models.Organization{
			ID:        uuid.New().String(),
			Name:      faker.Word() + " " + faker.Word() + " Inc.",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		orgs = append(orgs, org)
	}

	return orgs
}

// GenerateMockTeams generates a slice of mock Team models
func GenerateMockTeams(count int, orgID string) []models.Team {
	var teams []models.Team

	for i := 0; i < count; i++ {
		team := models.Team{
			ID:        uuid.New().String(),
			OrgID:     orgID,
			Name:      "Team " + faker.Word() + " " + faker.Word(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		teams = append(teams, team)
	}

	return teams
}

// GenerateMockTeamMembers generates mock team members
func GenerateMockTeamMembers(teamID string, userIDs []string) []models.TeamMember {
	var members []models.TeamMember
	roles := []string{"owner", "admin", "editor", "viewer"}

	for i, userID := range userIDs {
		member := models.TeamMember{
			TeamID:    teamID,
			UserID:    userID,
			Role:      roles[i%len(roles)],
			CreatedAt: time.Now(),
		}
		members = append(members, member)
	}

	return members
}

// GenerateMockSocialAccounts generates mock social media accounts
func GenerateMockSocialAccounts(count int, orgID string, teamID *string) []models.SocialAccount {
	var accounts []models.SocialAccount
	platforms := []string{"twitter", "facebook", "instagram", "linkedin", "tiktok"}

	for i := 0; i < count; i++ {
		account := models.SocialAccount{
			ID:          uuid.New().String(),
			OrgID:       orgID,
			TeamID:      teamID,
			Platform:    platforms[rand.Intn(len(platforms))],
			Handle:      "@" + strings.ToLower(faker.Username()),
			ExternalID:  strconv.Itoa(10000000 + rand.Intn(90000000)),
			AuthKind:    "oauth2",
			AuthMeta:    map[string]interface{}{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		accounts = append(accounts, account)
	}

	return accounts
}

// GenerateMockOAuthProviders generates mock OAuth provider configurations
func GenerateMockOAuthProviders() []models.OAuthProvider {
	providers := []string{"google", "github", "microsoft", "facebook", "linkedin"}
	var result []models.OAuthProvider

	for i, provider := range providers {
		result = append(result, models.OAuthProvider{
			ID:           uuid.New().String(),
			Provider:     provider,
			ClientID:     "client-" + uuid.New().String()[:8],
			ClientSecret: "secret-" + uuid.New().String()[:16],
			AuthURL:      "https://" + provider + ".com/oauth2/auth",
			TokenURL:     "https://" + provider + ".com/oauth2/token",
			UserInfoURL:  "https://api." + provider + ".com/userinfo",
			Scopes:       []string{"openid", "profile", "email"},
			IsActive:     i%2 == 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}
	return result
}

// GenerateMockOAuthAccounts generates mock OAuth account links
func GenerateMockOAuthAccounts(users []models.User, providers []models.OAuthProvider) []models.OAuthAccount {
	var accounts []models.OAuthAccount
	if len(providers) == 0 || len(users) == 0 {
		return accounts
	}

	for _, user := range users {
		// Link each user to 1-2 OAuth providers
		providerCount := 1 + rand.Intn(2)
		usedIndices := make(map[int]bool)

		for i := 0; i < providerCount && i < len(providers); i++ {
			// Find a provider index we haven't used yet
			var idx int
			for {
				idx = rand.Intn(len(providers))
				if !usedIndices[idx] {
					usedIndices[idx] = true
					break
				}
			}

			provider := providers[idx]
			accounts = append(accounts, models.OAuthAccount{
				ID:             uuid.New().String(),
				UserID:         user.ID,
				Provider:       provider.Provider,
				ProviderUserID: "oauth2|" + uuid.New().String()[:16],
				AccessToken:    "access-" + uuid.New().String(),
				RefreshToken:   "refresh-" + uuid.New().String(),
				TokenType:      "Bearer",
				ExpiresAt:      time.Now().Add(24 * time.Hour),
				Scope:          "openid profile email",
				Email:          user.Email,
				ProfileData:    map[string]interface{}{"name": user.DisplayName, "provider": provider.Provider},
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			})
		}
	}
	return accounts
}

// GenerateMockPhoneVerifications generates mock phone verifications
func GenerateMockPhoneVerifications(count int) []models.PhoneVerification {
	var verifications []models.PhoneVerification

	for i := 0; i < count; i++ {
		code := fmt.Sprintf("%06d", rand.Intn(1000000))
		expiresAt := time.Now().Add(10 * time.Minute)
		
		verifications = append(verifications, models.PhoneVerification{
			ID:          uuid.New().String(),
			PhoneNumber: "+1" + fmt.Sprintf("%010d", 2000000000+rand.Intn(799999999)),
			Code:        code, // In a real implementation, this would be hashed
			Attempts:    0,
			ExpiresAt:   expiresAt,
			Verified:    i%2 == 0,
			CreatedAt:   time.Now(),
			VerifiedAt:  time.Now().Add(-5 * time.Minute),
		})
	}
	return verifications
}

// GenerateMockSocialAccountLinks generates mock social account links
func GenerateMockSocialAccountLinks(users []models.User) []models.SocialAccountLink {
	var links []models.SocialAccountLink
	platforms := []string{"twitter", "facebook", "instagram", "linkedin", "tiktok"}

	for _, user := range users {
		// Each user gets 1-3 social account links
		count := 1 + rand.Intn(3)
		usedIndices := make(map[int]bool)

		for i := 0; i < count && i < len(platforms); i++ {
			// Find a platform index we haven't used yet
			var idx int
			for {
				idx = rand.Intn(len(platforms))
				if !usedIndices[idx] {
					usedIndices[idx] = true
					break
				}
			}

			platform := platforms[idx]
			username := strings.ToLower(strings.ReplaceAll(user.DisplayName, " ", ".")) + strconv.Itoa(rand.Intn(100))
			
			links = append(links, models.SocialAccountLink{
				ID:                uuid.New().String(),
				UserID:            user.ID,
				Provider:          platform,
				ExternalID:        strconv.Itoa(10000000 + rand.Intn(90000000)),
				Username:          username,
				DisplayName:       user.DisplayName + "'s " + strings.Title(platform),
				ProfileURL:        "https://" + platform + ".com/" + username,
				AvatarURL:         "https://api.dicebear.com/7.x/avataaars/svg?seed=" + username,
				AccessToken:       "access-" + uuid.New().String(),
				AccessTokenSecret: "secret-" + uuid.New().String()[:16],
				RefreshToken:      "refresh-" + uuid.New().String(),
				TokenExpiry:       timePtr(time.Now().Add(30 * 24 * time.Hour)),
				IsPrimary:         i == 0, // First one is primary
				Metadata: map[string]interface{}{
					"follower_count": 1000 + rand.Intn(10000),
					"verified":       i%3 == 0,
				},
				LastSyncedAt: timePtr(time.Now().Add(-time.Duration(rand.Intn(24)) * time.Hour)),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			})
		}
	}
	return links
}

// Helper function to get a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
