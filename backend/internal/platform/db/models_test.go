package db

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/your-org/ai-tools-atlas-backend/internal/auth"
	"github.com/your-org/ai-tools-atlas-backend/internal/badges"
	"github.com/your-org/ai-tools-atlas-backend/internal/bookmarks"
	"github.com/your-org/ai-tools-atlas-backend/internal/categories"
	"github.com/your-org/ai-tools-atlas-backend/internal/reviews"
	"github.com/your-org/ai-tools-atlas-backend/internal/tools"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping model tests")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Failed to get sql.DB: %v", err)
		return
	}

	// Drop all tables in reverse dependency order
	tables := []string{
		"bookmarks", "reviews", "users",
		"tool_alternatives", "tool_badges", "badges",
		"media", "tool_tags", "tags", "tools", "categories",
	}

	for _, table := range tables {
		sqlDB.Exec("DROP TABLE IF EXISTS " + table + " CASCADE")
	}
}

func TestGORMAutoMigration(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	t.Run("auto-migrate all models successfully", func(t *testing.T) {
		err := db.AutoMigrate(
			&categories.Category{},
			&tools.Tool{},
			&tools.Tag{},
			&tools.Media{},
			&badges.Badge{},
			&tools.ToolBadge{},
			&tools.ToolAlternative{},
			&auth.User{},
			&reviews.Review{},
			&bookmarks.Bookmark{},
		)

		if err != nil {
			t.Fatalf("AutoMigrate failed: %v", err)
		}

		// Verify tables exist
		sqlDB, _ := db.DB()
		tables := []string{
			"categories", "tools", "tags", "media", "badges",
			"tool_alternatives", "users", "reviews", "bookmarks",
		}

		for _, table := range tables {
			var exists bool
			query := `SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)`
			err := sqlDB.QueryRow(query, table).Scan(&exists)
			if err != nil {
				t.Fatalf("Failed to check table %s: %v", table, err)
			}
			if !exists {
				t.Errorf("Table %s was not created by AutoMigrate", table)
			}
		}
	})
}

func TestModelRelationships(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// AutoMigrate
	db.AutoMigrate(
		&categories.Category{},
		&tools.Tool{},
		&tools.Tag{},
		&tools.Media{},
		&badges.Badge{},
		&tools.ToolBadge{},
	)

	t.Run("tool has relationship with category", func(t *testing.T) {
		// Create category
		cat := categories.Category{
			Slug: "productivity",
			Name: "Productivity",
		}
		db.Create(&cat)

		// Create tool with category
		tool := tools.Tool{
			Slug:              "test-tool",
			Name:              "Test Tool",
			PrimaryCategoryID: cat.ID,
		}
		db.Create(&tool)

		// Fetch tool with preloaded category
		var fetchedTool tools.Tool
		db.Preload("PrimaryCategory").First(&fetchedTool, tool.ID)

		if fetchedTool.PrimaryCategory.Slug != "productivity" {
			t.Errorf("Expected category slug 'productivity', got '%s'", fetchedTool.PrimaryCategory.Slug)
		}
	})

	t.Run("tool has many-to-many relationship with tags", func(t *testing.T) {
		// Create tool
		tool := tools.Tool{
			Slug: "test-tool-2",
			Name: "Test Tool 2",
		}
		db.Create(&tool)

		// Create tags
		tag1 := tools.Tag{Slug: "ai", Name: "AI"}
		tag2 := tools.Tag{Slug: "ml", Name: "ML"}
		db.Create(&tag1)
		db.Create(&tag2)

		// Associate tags with tool
		db.Model(&tool).Association("Tags").Append(&tag1, &tag2)

		// Fetch tool with tags
		var fetchedTool tools.Tool
		db.Preload("Tags").First(&fetchedTool, tool.ID)

		if len(fetchedTool.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(fetchedTool.Tags))
		}
	})

	t.Run("tool has many media", func(t *testing.T) {
		// Create tool
		tool := tools.Tool{
			Slug: "test-tool-3",
			Name: "Test Tool 3",
		}
		db.Create(&tool)

		// Create media
		media1 := tools.Media{
			ToolID: tool.ID,
			Type:   "screenshot",
			URL:    "https://example.com/screenshot1.png",
		}
		media2 := tools.Media{
			ToolID: tool.ID,
			Type:   "video",
			URL:    "https://example.com/video1.mp4",
		}
		db.Create(&media1)
		db.Create(&media2)

		// Fetch tool with media
		var fetchedTool tools.Tool
		db.Preload("Media").First(&fetchedTool, tool.ID)

		if len(fetchedTool.Media) != 2 {
			t.Errorf("Expected 2 media items, got %d", len(fetchedTool.Media))
		}
	})

	t.Run("tool has many-to-many relationship with badges", func(t *testing.T) {
		// Create tool
		tool := tools.Tool{
			Slug: "test-tool-4",
			Name: "Test Tool 4",
		}
		db.Create(&tool)

		// Create badges
		badge1 := badges.Badge{Slug: "top-rated", Name: "Top Rated"}
		badge2 := badges.Badge{Slug: "editors-choice", Name: "Editor's Choice"}
		db.Create(&badge1)
		db.Create(&badge2)

		// Associate badges with tool
		db.Model(&tool).Association("Badges").Append(&badge1, &badge2)

		// Fetch tool with badges
		var fetchedTool tools.Tool
		db.Preload("Badges").First(&fetchedTool, tool.ID)

		if len(fetchedTool.Badges) != 2 {
			t.Errorf("Expected 2 badges, got %d", len(fetchedTool.Badges))
		}
	})
}

func TestJSONMarshaling(t *testing.T) {
	t.Run("tool marshals to snake_case JSON", func(t *testing.T) {
		tool := tools.Tool{
			ID:               1,
			Slug:             "test-tool",
			Name:             "Test Tool",
			LogoURL:          "https://example.com/logo.png",
			PrimaryCategoryID: 5,
			AvgRatingOverall: 4.5,
			ReviewCount:      100,
			BookmarkCount:    50,
			TrendingScore:    95.5,
			HasFreeTier:      true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		jsonData, err := json.Marshal(tool)
		if err != nil {
			t.Fatalf("Failed to marshal tool to JSON: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(jsonData, &result)

		// Verify snake_case keys exist
		expectedKeys := []string{
			"id", "slug", "name", "logo_url",
			"primary_category_id", "avg_rating_overall",
			"review_count", "bookmark_count", "trending_score",
			"has_free_tier", "created_at", "updated_at",
		}

		for _, key := range expectedKeys {
			if _, exists := result[key]; !exists {
				t.Errorf("Expected key '%s' not found in JSON", key)
			}
		}

		// Verify camelCase keys don't exist
		if _, exists := result["primaryCategoryId"]; exists {
			t.Error("Found camelCase key 'primaryCategoryId', expected snake_case")
		}
	})

	t.Run("category marshals to snake_case JSON", func(t *testing.T) {
		cat := categories.Category{
			ID:           1,
			Slug:         "productivity",
			Name:         "Productivity",
			IconURL:      "https://example.com/icon.png",
			DisplayOrder: 1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		jsonData, err := json.Marshal(cat)
		if err != nil {
			t.Fatalf("Failed to marshal category to JSON: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(jsonData, &result)

		if _, exists := result["icon_url"]; !exists {
			t.Error("Expected 'icon_url' key in JSON")
		}

		if _, exists := result["display_order"]; !exists {
			t.Error("Expected 'display_order' key in JSON")
		}
	})

	t.Run("user password_hash is not exposed in JSON", func(t *testing.T) {
		user := auth.User{
			ID:           1,
			Email:        "test@example.com",
			PasswordHash: "hashed_password_should_not_appear",
			DisplayName:  "Test User",
			Role:         "user",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal user to JSON: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(jsonData, &result)

		if _, exists := result["password_hash"]; exists {
			t.Error("password_hash should not be exposed in JSON")
		}

		if _, exists := result["PasswordHash"]; exists {
			t.Error("PasswordHash should not be exposed in JSON")
		}
	})

	t.Run("review marshals all rating fields correctly", func(t *testing.T) {
		easeOfUse := 5
		value := 4

		review := reviews.Review{
			ID:              1,
			ToolID:          10,
			UserID:          20,
			RatingOverall:   5,
			RatingEaseOfUse: &easeOfUse,
			RatingValue:     &value,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		jsonData, err := json.Marshal(review)
		if err != nil {
			t.Fatalf("Failed to marshal review to JSON: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(jsonData, &result)

		expectedKeys := []string{
			"id", "tool_id", "user_id", "rating_overall",
			"rating_ease_of_use", "rating_value",
			"created_at", "updated_at",
		}

		for _, key := range expectedKeys {
			if _, exists := result[key]; !exists {
				t.Errorf("Expected key '%s' not found in JSON", key)
			}
		}
	})
}

func TestGORMTags(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	db.AutoMigrate(
		&categories.Category{},
		&tools.Tool{},
		&auth.User{},
		&bookmarks.Bookmark{},
	)

	t.Run("unique constraints are enforced", func(t *testing.T) {
		// Create first category
		cat1 := categories.Category{
			Slug: "productivity",
			Name: "Productivity",
		}
		db.Create(&cat1)

		// Try to create duplicate slug
		cat2 := categories.Category{
			Slug: "productivity",
			Name: "Productivity 2",
		}
		result := db.Create(&cat2)

		if result.Error == nil {
			t.Error("Expected unique constraint violation, but insert succeeded")
		}
	})

	t.Run("composite unique constraint on bookmarks", func(t *testing.T) {
		// Create user and tool (simplified, actual FK constraints may require more setup)
		user := auth.User{Email: "user@example.com", PasswordHash: "hash", Role: "user"}
		db.Create(&user)

		tool := tools.Tool{Slug: "tool1", Name: "Tool 1"}
		db.Create(&tool)

		// Create first bookmark
		bm1 := bookmarks.Bookmark{
			UserID: user.ID,
			ToolID: tool.ID,
		}
		db.Create(&bm1)

		// Try to create duplicate bookmark
		bm2 := bookmarks.Bookmark{
			UserID: user.ID,
			ToolID: tool.ID,
		}
		result := db.Create(&bm2)

		if result.Error == nil {
			t.Error("Expected unique constraint violation on (user_id, tool_id), but insert succeeded")
		}
	})
}
