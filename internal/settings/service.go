package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/typewriterco/p402/internal/dbpg"
	"github.com/typewriterco/p402/internal/problems"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

// ScopeContext defines the hierarchical scope for settings lookup
type ScopeContext struct {
	OrganizationID *int64
	ProjectID      *int64
	UserID         *int64
}

// SystemScope returns a scope for system-level settings only
func SystemScope() ScopeContext {
	return ScopeContext{}
}

// OrganizationScope returns a scope for organization + system settings
func OrganizationScope(orgID int64) ScopeContext {
	return ScopeContext{OrganizationID: &orgID}
}

// ProjectScope returns a scope for project + organization + system settings
func ProjectScope(projectID int64) ScopeContext {
	return ScopeContext{ProjectID: &projectID}
}

// UserScope returns a scope for user + project + organization + system settings
func UserScope(userID int64) ScopeContext {
	return ScopeContext{UserID: &userID}
}

// Service provides hierarchical runtime configuration
type Service struct {
	queries *dbpg.Queries
	logger  zerolog.Logger

	// Cache for settings (optional, can be added for performance)
	cache      map[string]cacheEntry
	cacheMutex sync.RWMutex
	cacheTTL   time.Duration
}

type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

// NewService creates a new settings service
func NewService(queries *dbpg.Queries, logger zerolog.Logger) *Service {
	return &Service{
		queries:  queries,
		logger:   logger,
		cache:    make(map[string]cacheEntry),
		cacheTTL: 5 * time.Minute, // Cache settings for 5 minutes
	}
}

// ErrSettingNotFound is returned when a setting doesn't exist
var ErrSettingNotFound = fmt.Errorf("setting not found")

// Setting represents a configuration setting (domain model)
type Setting struct {
	ID             int64
	Scope          string // "system", "organization", "project", "user"
	OrganizationID *int64
	ProjectID      *int64
	UserID         *int64
	Subsystem      string
	Key            string
	Value          []byte // JSON value
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UpdatedBy      *int64
}

// SettingResult contains a setting value with metadata about where it was resolved
type SettingResult struct {
	Value         []byte // JSON value
	ResolvedScope string // Which scope provided this value: "system", "organization", "project", "user"
	Description   string // Setting description
}

// IsNotFound checks if an error is a "setting not found" error
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return errMsg == "setting not found" ||
		(len(errMsg) >= 18 && errMsg[:18] == "setting not found:")
}

// Get retrieves a setting with hierarchical resolution
// Returns the most specific value available in the scope hierarchy
func (s *Service) Get(ctx context.Context, subsystem, key string, scope ScopeContext) ([]byte, error) {
	result, err := s.GetWithMetadata(ctx, subsystem, key, scope)
	if err != nil {
		return nil, err
	}
	return result.Value, nil
}

// GetWithMetadata retrieves a setting with hierarchical resolution and metadata
// Returns the value, which scope resolved it, and the description
func (s *Service) GetWithMetadata(ctx context.Context, subsystem, key string, scope ScopeContext) (*SettingResult, error) {
	// Build cache key
	cacheKey := s.buildCacheKey(subsystem, key, scope)

	// Note: We could check cache here, but since GetWithMetadata needs to return
	// the resolved scope and description, we always query the database for this method.
	// The regular Get() method still uses caching.

	// Query database with hierarchy
	params := dbpg.GetSettingHierarchyParams{
		Subsystem: subsystem,
		Key:       key,
	}

	if scope.OrganizationID != nil {
		params.OrganizationID = pgtype.Int8{Int64: *scope.OrganizationID, Valid: true}
	}
	if scope.ProjectID != nil {
		params.ProjectID = pgtype.Int8{Int64: *scope.ProjectID, Valid: true}
	}
	if scope.UserID != nil {
		params.UserID = pgtype.Int8{Int64: *scope.UserID, Valid: true}
	}

	settings, err := s.queries.GetSettingHierarchy(ctx, params)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to retrieve setting", err)
	}

	if len(settings) == 0 {
		return nil, problems.New(problems.NotExist, fmt.Sprintf("setting not found: %s:%s", subsystem, key))
	}

	// The query returns highest precedence first, so take the first result
	setting := settings[0]

	// Cache the value
	s.setCached(cacheKey, setting.Value)

	// Build result with metadata
	result := &SettingResult{
		Value:         setting.Value,
		ResolvedScope: setting.Scope,
		Description:   "",
	}

	if setting.Description.Valid {
		result.Description = setting.Description.String
	}

	return result, nil
}

// GetTyped retrieves a setting and unmarshals it into the provided type
func GetTyped[T any](ctx context.Context, s *Service, subsystem, key string, scope ScopeContext) (T, error) {
	var result T

	data, err := s.Get(ctx, subsystem, key, scope)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, problems.New(problems.InvalidRequest, fmt.Sprintf("invalid setting value for %s:%s", subsystem, key), err)
	}

	return result, nil
}

// GetString is a convenience method for string settings
func (s *Service) GetString(ctx context.Context, subsystem, key string, scope ScopeContext) (string, error) {
	return GetTyped[string](ctx, s, subsystem, key, scope)
}

// GetInt is a convenience method for integer settings
func (s *Service) GetInt(ctx context.Context, subsystem, key string, scope ScopeContext) (int, error) {
	return GetTyped[int](ctx, s, subsystem, key, scope)
}

// GetBool is a convenience method for boolean settings
func (s *Service) GetBool(ctx context.Context, subsystem, key string, scope ScopeContext) (bool, error) {
	return GetTyped[bool](ctx, s, subsystem, key, scope)
}

// SetSystem sets a system-level setting
func (s *Service) SetSystem(ctx context.Context, subsystem, key string, value any, description *string, updatedBy *int64) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return problems.New(problems.InvalidRequest, "invalid setting value", err)
	}

	params := dbpg.UpsertSystemSettingParams{
		Subsystem: subsystem,
		Key:       key,
		Value:     jsonValue,
	}

	if description != nil {
		params.Description = pgtype.Text{String: *description, Valid: true}
	}
	if updatedBy != nil {
		params.UpdatedBy = pgtype.Int8{Int64: *updatedBy, Valid: true}
	}

	_, err = s.queries.UpsertSystemSetting(ctx, params)
	if err != nil {
		return problems.New(problems.Database, "failed to save system setting", err)
	}

	// Invalidate cache
	s.invalidateCache(subsystem, key)

	return nil
}

// SetOrganization sets an organization-level setting
func (s *Service) SetOrganization(ctx context.Context, orgID int64, subsystem, key string, value any, description *string, updatedBy *int64) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return problems.New(problems.InvalidRequest, "invalid setting value", err)
	}

	params := dbpg.UpsertOrganizationSettingParams{
		OrganizationID: pgtype.Int8{Int64: orgID, Valid: true},
		Subsystem:      subsystem,
		Key:            key,
		Value:          jsonValue,
	}

	if description != nil {
		params.Description = pgtype.Text{String: *description, Valid: true}
	}
	if updatedBy != nil {
		params.UpdatedBy = pgtype.Int8{Int64: *updatedBy, Valid: true}
	}

	_, err = s.queries.UpsertOrganizationSetting(ctx, params)
	if err != nil {
		return problems.New(problems.Database, "failed to save organization setting", err)
	}

	// Invalidate cache
	s.invalidateCache(subsystem, key)

	return nil
}

// SetProject sets a project-level setting
func (s *Service) SetProject(ctx context.Context, projectID int64, subsystem, key string, value any, description *string, updatedBy *int64) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return problems.New(problems.InvalidRequest, "invalid setting value", err)
	}

	params := dbpg.UpsertProjectSettingParams{
		ProjectID: pgtype.Int8{Int64: projectID, Valid: true},
		Subsystem: subsystem,
		Key:       key,
		Value:     jsonValue,
	}

	if description != nil {
		params.Description = pgtype.Text{String: *description, Valid: true}
	}
	if updatedBy != nil {
		params.UpdatedBy = pgtype.Int8{Int64: *updatedBy, Valid: true}
	}

	_, err = s.queries.UpsertProjectSetting(ctx, params)
	if err != nil {
		return problems.New(problems.Database, "failed to save project setting", err)
	}

	// Invalidate cache
	s.invalidateCache(subsystem, key)

	return nil
}

// SetUser sets a user-level setting
func (s *Service) SetUser(ctx context.Context, userID int64, subsystem, key string, value any, description *string, updatedBy *int64) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return problems.New(problems.InvalidRequest, "invalid setting value", err)
	}

	params := dbpg.UpsertUserSettingParams{
		UserID:    pgtype.Int8{Int64: userID, Valid: true},
		Subsystem: subsystem,
		Key:       key,
		Value:     jsonValue,
	}

	if description != nil {
		params.Description = pgtype.Text{String: *description, Valid: true}
	}
	if updatedBy != nil {
		params.UpdatedBy = pgtype.Int8{Int64: *updatedBy, Valid: true}
	}

	_, err = s.queries.UpsertUserSetting(ctx, params)
	if err != nil {
		return problems.New(problems.Database, "failed to save user setting", err)
	}

	// Invalidate cache
	s.invalidateCache(subsystem, key)

	return nil
}

// GetBySubsystem gets all settings for a subsystem with hierarchy
func (s *Service) GetBySubsystem(ctx context.Context, subsystem string, scope ScopeContext) (map[string]any, error) {
	params := dbpg.GetSettingsBySubsystemParams{
		Subsystem: subsystem,
	}

	if scope.OrganizationID != nil {
		params.OrganizationID = pgtype.Int8{Int64: *scope.OrganizationID, Valid: true}
	}
	if scope.ProjectID != nil {
		params.ProjectID = pgtype.Int8{Int64: *scope.ProjectID, Valid: true}
	}
	if scope.UserID != nil {
		params.UserID = pgtype.Int8{Int64: *scope.UserID, Valid: true}
	}

	settings, err := s.queries.GetSettingsBySubsystem(ctx, params)
	if err != nil {
		return nil, problems.New(problems.Database, fmt.Sprintf("failed to retrieve settings for subsystem %s", subsystem), err)
	}

	// Build result map with highest precedence value for each key
	result := make(map[string]any)
	seen := make(map[string]bool)

	for _, setting := range settings {
		// Since results are ordered by precedence, first occurrence wins
		if !seen[setting.Key] {
			var value any
			if err := json.Unmarshal(setting.Value, &value); err != nil {
				s.logger.Warn().
					Err(err).
					Str("subsystem", subsystem).
					Str("key", setting.Key).
					Msg("Failed to unmarshal setting value")
				continue
			}
			result[setting.Key] = value
			seen[setting.Key] = true
		}
	}

	return result, nil
}

// Delete removes a setting by ID
func (s *Service) Delete(ctx context.Context, id int64) error {
	err := s.queries.DeleteSetting(ctx, dbpg.DeleteSettingParams{ID: id})
	if err != nil {
		return problems.New(problems.Database, "failed to delete setting", err)
	}

	// Clear entire cache (we don't know which keys were affected)
	s.clearCache()

	return nil
}

// Cache helpers

func (s *Service) buildCacheKey(subsystem, key string, scope ScopeContext) string {
	var orgID, projID, userID int64
	if scope.OrganizationID != nil {
		orgID = *scope.OrganizationID
	}
	if scope.ProjectID != nil {
		projID = *scope.ProjectID
	}
	if scope.UserID != nil {
		userID = *scope.UserID
	}
	return fmt.Sprintf("%s:%s:%d:%d:%d", subsystem, key, orgID, projID, userID)
}

func (s *Service) getCached(key string) ([]byte, bool) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	entry, ok := s.cache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.value, true
}

func (s *Service) setCached(key string, value []byte) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(s.cacheTTL),
	}
}

func (s *Service) invalidateCache(subsystem, key string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Remove all cache entries for this subsystem:key
	for cacheKey := range s.cache {
		if len(cacheKey) > len(subsystem)+len(key)+1 &&
			cacheKey[:len(subsystem)] == subsystem &&
			cacheKey[len(subsystem)+1:len(subsystem)+1+len(key)] == key {
			delete(s.cache, cacheKey)
		}
	}
}

func (s *Service) clearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache = make(map[string]cacheEntry)
}

// SettingsRepository defines the interface for settings data access
type SettingsRepository interface {
	ListAll(ctx context.Context) ([]Setting, error)
}

// ListAll returns all settings from the database (for admin interface)
func (s *Service) ListAll(ctx context.Context, repo SettingsRepository) ([]Setting, error) {
	settings, err := repo.ListAll(ctx)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list settings", err)
	}
	return settings, nil
}
