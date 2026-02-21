package settings

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/richardbowden/degrees/internal/dbpg"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper to setup test database and settings service
func setupTestService(t *testing.T) (*Service, *pgxpool.Pool, func()) {
	t.Helper()

	// Get database URL from environment or use default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://p402:p402@localhost:5432/p402_test?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil, nil, func() {}
	}

	// Verify database is actually accessible
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("Skipping test: cannot ping test database: %v", err)
		return nil, nil, func() {}
	}

	// Verify settings table exists
	var exists bool
	err = pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'settings')").Scan(&exists)
	if err != nil || !exists {
		pool.Close()
		t.Skipf("Skipping test: settings table does not exist. Run migrations first.")
		return nil, nil, func() {}
	}

	// Create queries
	queries := dbpg.New(pool)

	// Create service
	logger := zerolog.Nop() // No-op logger for tests
	service := NewService(queries, logger)

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		_, err := pool.Exec(context.Background(), "DELETE FROM settings WHERE subsystem LIKE 'test%'")
		if err != nil {
			t.Logf("cleanup error: %v", err)
		}
		pool.Close()
	}

	return service, pool, cleanup
}

func TestHierarchicalResolution(t *testing.T) {
	service, _, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Setup: Create settings at different scopes
	orgID := int64(123)
	projectID := int64(456)

	// System default
	err := service.SetSystem(ctx, "test_features", "feature_a", true, strPtr("System default"), nil)
	require.NoError(t, err)

	// Organization override
	err = service.SetOrganization(ctx, orgID, "test_features", "feature_a", false, strPtr("Org override"), nil)
	require.NoError(t, err)

	// Project override
	err = service.SetProject(ctx, projectID, "test_features", "feature_a", true, strPtr("Project override"), nil)
	require.NoError(t, err)

	// Test: System scope gets system default
	t.Run("system scope", func(t *testing.T) {
		value, err := service.GetBool(ctx, "test_features", "feature_a", SystemScope())
		require.NoError(t, err)
		assert.True(t, value, "system scope should return system default (true)")
	})

	// Test: Organization scope gets org override
	t.Run("organization scope", func(t *testing.T) {
		value, err := service.GetBool(ctx, "test_features", "feature_a", OrganizationScope(orgID))
		require.NoError(t, err)
		assert.False(t, value, "org scope should return org override (false)")
	})

	// Test: Project scope gets project override
	t.Run("project scope", func(t *testing.T) {
		value, err := service.GetBool(ctx, "test_features", "feature_a", ProjectScope(projectID))
		require.NoError(t, err)
		assert.True(t, value, "project scope should return project override (true)")
	})

	// Test: Different org gets system default
	t.Run("different organization", func(t *testing.T) {
		value, err := service.GetBool(ctx, "test_features", "feature_a", OrganizationScope(999))
		require.NoError(t, err)
		assert.True(t, value, "different org should get system default (true)")
	})
}

func TestTypedGetters(t *testing.T) {
	service, _, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Test GetBool
	t.Run("GetBool", func(t *testing.T) {
		err := service.SetSystem(ctx, "test_types", "bool_setting", true, nil, nil)
		require.NoError(t, err)

		value, err := service.GetBool(ctx, "test_types", "bool_setting", SystemScope())
		require.NoError(t, err)
		assert.True(t, value)
	})

	// Test GetInt
	t.Run("GetInt", func(t *testing.T) {
		err := service.SetSystem(ctx, "test_types", "int_setting", 42, nil, nil)
		require.NoError(t, err)

		value, err := service.GetInt(ctx, "test_types", "int_setting", SystemScope())
		require.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	// Test GetString
	t.Run("GetString", func(t *testing.T) {
		err := service.SetSystem(ctx, "test_types", "string_setting", "hello", nil, nil)
		require.NoError(t, err)

		value, err := service.GetString(ctx, "test_types", "string_setting", SystemScope())
		require.NoError(t, err)
		assert.Equal(t, "hello", value)
	})

	// Test GetTyped with struct
	t.Run("GetTyped struct", func(t *testing.T) {
		type RateLimitConfig struct {
			Limit  int    `json:"limit"`
			Window string `json:"window"`
		}

		config := RateLimitConfig{
			Limit:  1000,
			Window: "1h",
		}

		err := service.SetSystem(ctx, "test_types", "rate_limit", config, nil, nil)
		require.NoError(t, err)

		value, err := GetTyped[RateLimitConfig](ctx, service, "test_types", "rate_limit", SystemScope())
		require.NoError(t, err)
		assert.Equal(t, 1000, value.Limit)
		assert.Equal(t, "1h", value.Window)
	})
}

func TestGetBySubsystem(t *testing.T) {
	service, _, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()
	orgID := int64(100)

	// Setup: Create multiple settings in a subsystem
	err := service.SetSystem(ctx, "test_limits", "max_users", 10, nil, nil)
	require.NoError(t, err)

	err = service.SetSystem(ctx, "test_limits", "max_projects", 5, nil, nil)
	require.NoError(t, err)

	err = service.SetOrganization(ctx, orgID, "test_limits", "max_users", 100, nil, nil)
	require.NoError(t, err)

	// Test: Get all settings for system scope
	t.Run("system scope", func(t *testing.T) {
		settings, err := service.GetBySubsystem(ctx, "test_limits", SystemScope())
		require.NoError(t, err)

		assert.Len(t, settings, 2)
		assert.Equal(t, float64(10), settings["max_users"]) // JSON unmarshals numbers as float64
		assert.Equal(t, float64(5), settings["max_projects"])
	})

	// Test: Get all settings for org scope (includes overrides)
	t.Run("organization scope", func(t *testing.T) {
		settings, err := service.GetBySubsystem(ctx, "test_limits", OrganizationScope(orgID))
		require.NoError(t, err)

		assert.Len(t, settings, 2)
		assert.Equal(t, float64(100), settings["max_users"]) // Org override
		assert.Equal(t, float64(5), settings["max_projects"]) // System default
	})
}

func TestCaching(t *testing.T) {
	service, pool, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Set a value
	err := service.SetSystem(ctx, "test_cache", "cached_value", "initial", nil, nil)
	require.NoError(t, err)

	// Get value (should cache it)
	value1, err := service.GetString(ctx, "test_cache", "cached_value", SystemScope())
	require.NoError(t, err)
	assert.Equal(t, "initial", value1)

	// Update value directly in database (bypass service to test cache)
	_, err = pool.Exec(ctx,
		"UPDATE settings SET value = '\"updated\"'::jsonb WHERE subsystem = 'test_cache' AND key = 'cached_value'")
	require.NoError(t, err)

	// Get value again immediately (should return cached value)
	value2, err := service.GetString(ctx, "test_cache", "cached_value", SystemScope())
	require.NoError(t, err)
	assert.Equal(t, "initial", value2, "should return cached value")

	// Update through service (should invalidate cache)
	err = service.SetSystem(ctx, "test_cache", "cached_value", "new", nil, nil)
	require.NoError(t, err)

	// Get value (should return new value)
	value3, err := service.GetString(ctx, "test_cache", "cached_value", SystemScope())
	require.NoError(t, err)
	assert.Equal(t, "new", value3, "cache should be invalidated after update")
}

func TestCacheExpiration(t *testing.T) {
	service, _, cleanup := setupTestService(t)
	defer cleanup()

	// Set short cache TTL for testing
	service.cacheTTL = 100 * time.Millisecond

	ctx := context.Background()

	// Set a value
	err := service.SetSystem(ctx, "test_expire", "value", "cached", nil, nil)
	require.NoError(t, err)

	// Get value (caches it)
	value1, err := service.GetString(ctx, "test_expire", "value", SystemScope())
	require.NoError(t, err)
	assert.Equal(t, "cached", value1)

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Update value
	err = service.SetSystem(ctx, "test_expire", "value", "expired", nil, nil)
	require.NoError(t, err)

	// Get value (should fetch from DB since cache expired)
	value2, err := service.GetString(ctx, "test_expire", "value", SystemScope())
	require.NoError(t, err)
	assert.Equal(t, "expired", value2)
}

func TestDelete(t *testing.T) {
	service, pool, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Create a setting
	err := service.SetSystem(ctx, "test_delete", "to_delete", "value", nil, nil)
	require.NoError(t, err)

	// Verify it exists
	_, err = service.GetString(ctx, "test_delete", "to_delete", SystemScope())
	require.NoError(t, err)

	// Get the setting ID by querying directly
	var settingID int64
	err = pool.QueryRow(ctx,
		"SELECT id FROM settings WHERE subsystem = 'test_delete' AND key = 'to_delete' AND scope = 'system'",
	).Scan(&settingID)
	require.NoError(t, err, "should find the setting")

	// Delete it
	err = service.Delete(ctx, settingID)
	require.NoError(t, err)

	// Verify it's gone
	_, err = service.GetString(ctx, "test_delete", "to_delete", SystemScope())
	require.Error(t, err, "setting should not be found after deletion")
}

func TestSettingNotFound(t *testing.T) {
	service, _, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Try to get a non-existent setting
	_, err := service.GetString(ctx, "test_nonexistent", "does_not_exist", SystemScope())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "setting not found")
}

func TestUpdateDescription(t *testing.T) {
	service, pool, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Create setting with description
	desc1 := "Original description"
	err := service.SetSystem(ctx, "test_desc", "setting", "value", &desc1, nil)
	require.NoError(t, err)

	// Update value, keeping description
	err = service.SetSystem(ctx, "test_desc", "setting", "new_value", nil, nil)
	require.NoError(t, err)

	// Verify description was kept (check in database)
	var description string
	err = pool.QueryRow(ctx,
		"SELECT description FROM settings WHERE subsystem = 'test_desc' AND key = 'setting'",
	).Scan(&description)
	require.NoError(t, err)
	assert.Equal(t, desc1, description)

	// Update with new description
	desc2 := "New description"
	err = service.SetSystem(ctx, "test_desc", "setting", "newer_value", &desc2, nil)
	require.NoError(t, err)

	// Verify description was updated
	err = pool.QueryRow(ctx,
		"SELECT description FROM settings WHERE subsystem = 'test_desc' AND key = 'setting'",
	).Scan(&description)
	require.NoError(t, err)
	assert.Equal(t, desc2, description)
}

func TestAuditTrail(t *testing.T) {
	service, pool, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()
	userID := int64(42)

	// Create setting with user ID
	err := service.SetSystem(ctx, "test_audit", "setting", "value", nil, &userID)
	require.NoError(t, err)

	// Verify updated_by was set
	var updatedBy *int64
	err = pool.QueryRow(ctx,
		"SELECT updated_by FROM settings WHERE subsystem = 'test_audit' AND key = 'setting'",
	).Scan(&updatedBy)
	require.NoError(t, err)
	require.NotNil(t, updatedBy)
	assert.Equal(t, userID, *updatedBy)
}

func TestMultipleScopesForSameKey(t *testing.T) {
	service, _, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Create settings at multiple scopes for the same key
	err := service.SetSystem(ctx, "test_multi", "color", "blue", nil, nil)
	require.NoError(t, err)

	err = service.SetOrganization(ctx, 100, "test_multi", "color", "red", nil, nil)
	require.NoError(t, err)

	err = service.SetOrganization(ctx, 200, "test_multi", "color", "green", nil, nil)
	require.NoError(t, err)

	err = service.SetProject(ctx, 300, "test_multi", "color", "yellow", nil, nil)
	require.NoError(t, err)

	// Test each scope gets the right value
	t.Run("system", func(t *testing.T) {
		value, err := service.GetString(ctx, "test_multi", "color", SystemScope())
		require.NoError(t, err)
		assert.Equal(t, "blue", value)
	})

	t.Run("org 100", func(t *testing.T) {
		value, err := service.GetString(ctx, "test_multi", "color", OrganizationScope(100))
		require.NoError(t, err)
		assert.Equal(t, "red", value)
	})

	t.Run("org 200", func(t *testing.T) {
		value, err := service.GetString(ctx, "test_multi", "color", OrganizationScope(200))
		require.NoError(t, err)
		assert.Equal(t, "green", value)
	})

	t.Run("project 300", func(t *testing.T) {
		value, err := service.GetString(ctx, "test_multi", "color", ProjectScope(300))
		require.NoError(t, err)
		assert.Equal(t, "yellow", value)
	})

	t.Run("org 999 (no override)", func(t *testing.T) {
		value, err := service.GetString(ctx, "test_multi", "color", OrganizationScope(999))
		require.NoError(t, err)
		assert.Equal(t, "blue", value, "should fall back to system default")
	})
}

func TestComplexJSONValues(t *testing.T) {
	service, _, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("nested struct", func(t *testing.T) {
		type EmailConfig struct {
			SMTP struct {
				Host     string `json:"host"`
				Port     int    `json:"port"`
				Username string `json:"username"`
			} `json:"smtp"`
			From struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"from"`
		}

		config := EmailConfig{}
		config.SMTP.Host = "smtp.example.com"
		config.SMTP.Port = 587
		config.SMTP.Username = "user@example.com"
		config.From.Name = "P402"
		config.From.Email = "noreply@p402.com"

		err := service.SetSystem(ctx, "test_json", "email_config", config, nil, nil)
		require.NoError(t, err)

		retrieved, err := GetTyped[EmailConfig](ctx, service, "test_json", "email_config", SystemScope())
		require.NoError(t, err)
		assert.Equal(t, "smtp.example.com", retrieved.SMTP.Host)
		assert.Equal(t, 587, retrieved.SMTP.Port)
		assert.Equal(t, "P402", retrieved.From.Name)
	})

	t.Run("array", func(t *testing.T) {
		roles := []string{"admin", "editor", "viewer"}

		err := service.SetSystem(ctx, "test_json", "allowed_roles", roles, nil, nil)
		require.NoError(t, err)

		retrieved, err := GetTyped[[]string](ctx, service, "test_json", "allowed_roles", SystemScope())
		require.NoError(t, err)
		assert.Equal(t, roles, retrieved)
	})

	t.Run("map", func(t *testing.T) {
		permissions := map[string]bool{
			"read":   true,
			"write":  true,
			"delete": false,
		}

		err := service.SetSystem(ctx, "test_json", "permissions", permissions, nil, nil)
		require.NoError(t, err)

		retrieved, err := GetTyped[map[string]bool](ctx, service, "test_json", "permissions", SystemScope())
		require.NoError(t, err)
		assert.Equal(t, permissions, retrieved)
	})
}

// Helper function
func strPtr(s string) *string {
	return &s
}
