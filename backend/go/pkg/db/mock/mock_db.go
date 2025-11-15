package mock

import (
	"context"
	"errors"

	"github.com/bitesinbyte/ferret/pkg/models"
	"gorm.io/gorm"
)

// MockDB implements a mock database for testing
type MockDB struct {
	users         map[string]models.User
	organizations map[string]models.Organization
	teams        map[string]models.Team
	content      map[string]interface{} // Using interface{} for flexibility
}

// NewMockDB creates a new mock database with sample data
func NewMockDB() *MockDB {
	db := &MockDB{
		users:         make(map[string]models.User),
		organizations: make(map[string]models.Organization),
		teams:        make(map[string]models.Team),
		content:      make(map[string]interface{}),
	}

	// Initialize with some test data
	orgs := GenerateMockOrganizations(2)
	for _, org := range orgs {
		db.organizations[org.ID] = org
		
		// Add users
		users := GenerateMockUsers(3, org.ID)
		for _, user := range users {
			db.users[user.ID] = user
		}
		
		// Add teams
		teams := GenerateMockTeams(2, org.ID)
		for _, team := range teams {
			db.teams[team.ID] = team
		}
	}
	
	return db
}

// User operations
func (m *MockDB) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		return &user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockDB) CreateUser(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	m.users[user.ID] = *user
	return nil
}

// Organization operations
func (m *MockDB) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	if org, exists := m.organizations[id]; exists {
		return &org, nil
	}
	return nil, gorm.ErrRecordNotFound
}

// Team operations
func (m *MockDB) GetTeamByID(ctx context.Context, id string) (*models.Team, error) {
	if team, exists := m.teams[id]; exists {
		return &team, nil
	}
	return nil, gorm.ErrRecordNotFound
}

// MockTx implements a mock database transaction
type MockTx struct {
	db *MockDB
}

func (m *MockDB) Begin() *MockTx {
	return &MockTx{db: m}
}

func (tx *MockTx) Commit() error {
	return nil
}

func (tx *MockTx) Rollback() error {
	return nil
}

// Test helpers
func (m *MockDB) AddTestUser(user models.User) {
	m.users[user.ID] = user
}

func (m *MockDB) AddTestOrganization(org models.Organization) {
	m.organizations[org.ID] = org
}

func (m *MockDB) AddTestTeam(team models.Team) {
	m.teams[team.ID] = team
}

// ClearTestData resets all test data
func (m *MockDB) ClearTestData() {
	m.users = make(map[string]models.User)
	m.organizations = make(map[string]models.Organization)
	m.teams = make(map[string]models.Team)
	m.content = make(map[string]interface{})
}
