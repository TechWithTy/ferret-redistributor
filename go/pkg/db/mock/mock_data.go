package mock

import (
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
	platforms := []string{"instagram", "linkedin", "twitter", "facebook", "youtube", "behiiv"}

	for i := 0; i < count && i < len(platforms); i++ {
		account := models.SocialAccount{
			ID:          uuid.New().String(),
			OrgID:       orgID,
			TeamID:      teamID,
			Platform:    platforms[i],
			Handle:      "@" + faker.Username(),
			ExternalID:  faker.UUIDDigit(),
			AuthKind:    "oauth2",
			AuthMeta:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		accounts = append(accounts, account)
	}

	return accounts
}
