package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestMigrations(t *testing.T) {
	// Skip if no test database configured
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping migration tests")
	}

	// Connect to database
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatalf("Failed to create postgres driver: %v", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../migrations",
		"postgres", driver)
	if err != nil {
		t.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Test: Migration Up
	t.Run("migrate up creates all tables", func(t *testing.T) {
		// Drop all tables first (clean state)
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			t.Logf("Warning: Could not migrate down: %v", err)
		}

		// Run migrations up
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			t.Fatalf("Migration up failed: %v", err)
		}

		// Verify all tables exist
		tables := []string{
			"categories", "tools", "tags", "tool_tags",
			"media", "badges", "tool_badges", "tool_alternatives",
			"users", "reviews", "bookmarks",
		}

		for _, table := range tables {
			var exists bool
			query := `SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)`
			err := db.QueryRow(query, table).Scan(&exists)
			if err != nil {
				t.Fatalf("Failed to check if table %s exists: %v", table, err)
			}
			if !exists {
				t.Errorf("Table %s does not exist after migration up", table)
			}
		}
	})

	// Test: Indexes Exist
	t.Run("migration creates all indexes", func(t *testing.T) {
		indexes := []string{
			"idx_tools_slug",
			"idx_tools_primary_category_id",
			"idx_tools_avg_rating",
			"idx_tools_bookmark_count",
			"idx_tools_trending_score",
			"idx_categories_slug",
			"idx_reviews_tool_id",
			"idx_reviews_user_id",
			"idx_bookmarks_user_tool",
		}

		for _, index := range indexes {
			var exists bool
			query := `SELECT EXISTS (
				SELECT FROM pg_indexes
				WHERE schemaname = 'public'
				AND indexname = $1
			)`
			err := db.QueryRow(query, index).Scan(&exists)
			if err != nil {
				t.Fatalf("Failed to check if index %s exists: %v", index, err)
			}
			if !exists {
				t.Errorf("Index %s does not exist after migration up", index)
			}
		}
	})

	// Test: Foreign Keys
	t.Run("foreign key constraints are enforced", func(t *testing.T) {
		// Try to insert a tool with invalid category_id
		_, err := db.Exec(`
			INSERT INTO tools (slug, name, primary_category_id)
			VALUES ('test-tool', 'Test Tool', 99999)
		`)
		if err == nil {
			t.Error("Expected foreign key violation, but insert succeeded")
		}
	})

	// Test: CHECK Constraints
	t.Run("check constraints are enforced", func(t *testing.T) {
		// Create a valid category first for FK constraint
		db.Exec(`INSERT INTO categories (slug, name) VALUES ('test-cat', 'Test Category')`)
		defer db.Exec(`DELETE FROM categories WHERE slug = 'test-cat'`)

		var categoryID int
		db.QueryRow(`SELECT id FROM categories WHERE slug = 'test-cat'`).Scan(&categoryID)

		// Test media.type CHECK constraint
		_, err := db.Exec(`
			INSERT INTO tools (slug, name, primary_category_id) VALUES ('test-tool-1', 'Test Tool 1', $1)
		`, categoryID)
		if err != nil {
			t.Logf("Tool insert: %v", err)
		}
		var toolID int
		db.QueryRow(`SELECT id FROM tools WHERE slug = 'test-tool-1'`).Scan(&toolID)

		_, err = db.Exec(`
			INSERT INTO media (tool_id, type, url)
			VALUES ($1, 'invalid_type', 'http://example.com/image.png')
		`, toolID)
		if err == nil {
			t.Error("Expected CHECK constraint violation for media.type, but insert succeeded")
		}

		db.Exec(`DELETE FROM tools WHERE slug = 'test-tool-1'`)

		// Test tool_alternatives.relationship_type CHECK constraint
		db.Exec(`INSERT INTO tools (slug, name, primary_category_id) VALUES ('tool-a', 'Tool A', $1)`, categoryID)
		db.Exec(`INSERT INTO tools (slug, name, primary_category_id) VALUES ('tool-b', 'Tool B', $1)`, categoryID)
		var toolAID, toolBID int
		db.QueryRow(`SELECT id FROM tools WHERE slug = 'tool-a'`).Scan(&toolAID)
		db.QueryRow(`SELECT id FROM tools WHERE slug = 'tool-b'`).Scan(&toolBID)

		_, err = db.Exec(`
			INSERT INTO tool_alternatives (tool_id, alternative_tool_id, relationship_type)
			VALUES ($1, $2, 'invalid_relationship')
		`, toolAID, toolBID)
		if err == nil {
			t.Error("Expected CHECK constraint violation for tool_alternatives.relationship_type, but insert succeeded")
		}

		db.Exec(`DELETE FROM tools WHERE slug IN ('tool-a', 'tool-b')`)

		// Test users.role CHECK constraint
		_, err = db.Exec(`
			INSERT INTO users (email, password_hash, role)
			VALUES ('test@example.com', 'hash123', 'superadmin')
		`)
		if err == nil {
			t.Error("Expected CHECK constraint violation for users.role, but insert succeeded")
		}

		// Test reviews.rating_overall CHECK constraint (must be 1-5)
		db.Exec(`INSERT INTO users (email, password_hash, role) VALUES ('reviewer@example.com', 'hash123', 'user')`)
		db.Exec(`INSERT INTO tools (slug, name, primary_category_id) VALUES ('test-tool-2', 'Test Tool 2', $1)`, categoryID)
		var userID, toolID2 int
		db.QueryRow(`SELECT id FROM users WHERE email = 'reviewer@example.com'`).Scan(&userID)
		db.QueryRow(`SELECT id FROM tools WHERE slug = 'test-tool-2'`).Scan(&toolID2)

		_, err = db.Exec(`
			INSERT INTO reviews (tool_id, user_id, rating_overall)
			VALUES ($1, $2, 10)
		`, toolID2, userID)
		if err == nil {
			t.Error("Expected CHECK constraint violation for reviews.rating_overall (>5), but insert succeeded")
		}

		_, err = db.Exec(`
			INSERT INTO reviews (tool_id, user_id, rating_overall)
			VALUES ($1, $2, 0)
		`, toolID2, userID)
		if err == nil {
			t.Error("Expected CHECK constraint violation for reviews.rating_overall (<1), but insert succeeded")
		}

		// Test reviews.moderation_status CHECK constraint
		_, err = db.Exec(`
			INSERT INTO reviews (tool_id, user_id, rating_overall, moderation_status)
			VALUES ($1, $2, 5, 'invalid_status')
		`, toolID2, userID)
		if err == nil {
			t.Error("Expected CHECK constraint violation for reviews.moderation_status, but insert succeeded")
		}

		// Cleanup
		db.Exec(`DELETE FROM reviews WHERE user_id = $1`, userID)
		db.Exec(`DELETE FROM tools WHERE slug = 'test-tool-2'`)
		db.Exec(`DELETE FROM users WHERE email = 'reviewer@example.com'`)
	})

	// Test: Migration Down
	t.Run("migrate down drops all tables", func(t *testing.T) {
		// Run migrations down
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			t.Fatalf("Migration down failed: %v", err)
		}

		// Verify all tables are dropped
		tables := []string{
			"categories", "tools", "tags", "tool_tags",
			"media", "badges", "tool_badges", "tool_alternatives",
			"users", "reviews", "bookmarks",
		}

		for _, table := range tables {
			var exists bool
			query := `SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)`
			err := db.QueryRow(query, table).Scan(&exists)
			if err != nil {
				t.Fatalf("Failed to check if table %s exists: %v", table, err)
			}
			if exists {
				t.Errorf("Table %s still exists after migration down", table)
			}
		}
	})

	// Cleanup: Ensure we're in clean state for next test run
	m.Down()
}

func TestMigrationUpDown(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping migration tests")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatalf("Failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../migrations",
		"postgres", driver)
	if err != nil {
		t.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Clean state
	m.Down()

	// Test multiple up/down cycles
	t.Run("multiple migration cycles work correctly", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			// Up
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				t.Fatalf("Cycle %d: Migration up failed: %v", i, err)
			}

			// Verify a sample table exists
			var exists bool
			err := db.QueryRow(`SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = 'tools'
			)`).Scan(&exists)
			if err != nil || !exists {
				t.Fatalf("Cycle %d: tools table should exist after up", i)
			}

			// Down
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				t.Fatalf("Cycle %d: Migration down failed: %v", i, err)
			}

			// Verify table is gone
			err = db.QueryRow(`SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = 'tools'
			)`).Scan(&exists)
			if err != nil {
				t.Fatalf("Cycle %d: Query failed: %v", i, err)
			}
			if exists {
				t.Fatalf("Cycle %d: tools table should not exist after down", i)
			}
		}
	})

	// Cleanup
	m.Down()
}

// Comprehensive test validating schema details for all tables
func TestSchemaDetails(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping migration tests")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatalf("Failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../migrations",
		"postgres", driver)
	if err != nil {
		t.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Setup
	m.Down()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migration up failed: %v", err)
	}
	defer m.Down()

	// Helper function to verify columns exist
	verifyColumns := func(t *testing.T, tableName string, requiredColumns []string) {
		for _, colName := range requiredColumns {
			var exists bool
			query := `
				SELECT EXISTS (
					SELECT 1
					FROM information_schema.columns
					WHERE table_schema = 'public'
					AND table_name = $1
					AND column_name = $2
				)
			`
			err := db.QueryRow(query, tableName, colName).Scan(&exists)
			if err != nil {
				t.Fatalf("Failed to check column %s in %s: %v", colName, tableName, err)
			}
			if !exists {
				t.Errorf("Column %s does not exist in %s table", colName, tableName)
			}
		}
	}

	t.Run("categories table has correct columns", func(t *testing.T) {
		verifyColumns(t, "categories", []string{
			"id", "slug", "name", "description", "icon_url", "display_order", "created_at", "updated_at",
		})
	})

	t.Run("tools table has correct columns", func(t *testing.T) {
		verifyColumns(t, "tools", []string{
			"id", "slug", "name", "logo_url", "tagline", "description", "best_for",
			"primary_use_cases", "pricing_summary", "target_roles", "platforms",
			"has_free_tier", "official_url", "primary_category_id", "avg_rating_overall",
			"review_count", "bookmark_count", "trending_score", "created_at", "updated_at", "archived_at",
		})
	})

	t.Run("tags table has correct columns", func(t *testing.T) {
		verifyColumns(t, "tags", []string{
			"id", "slug", "name", "created_at",
		})
	})

	t.Run("tool_tags table has correct columns", func(t *testing.T) {
		verifyColumns(t, "tool_tags", []string{
			"tool_id", "tag_id",
		})
	})

	t.Run("media table has correct columns", func(t *testing.T) {
		verifyColumns(t, "media", []string{
			"id", "tool_id", "type", "url", "thumbnail_url", "display_order", "created_at",
		})
	})

	t.Run("badges table has correct columns", func(t *testing.T) {
		verifyColumns(t, "badges", []string{
			"id", "slug", "name", "description", "icon_url", "created_at",
		})
	})

	t.Run("tool_badges table has correct columns", func(t *testing.T) {
		verifyColumns(t, "tool_badges", []string{
			"tool_id", "badge_id", "assigned_at",
		})
	})

	t.Run("tool_alternatives table has correct columns", func(t *testing.T) {
		verifyColumns(t, "tool_alternatives", []string{
			"id", "tool_id", "alternative_tool_id", "relationship_type", "created_at",
		})
	})

	t.Run("users table has correct columns", func(t *testing.T) {
		verifyColumns(t, "users", []string{
			"id", "email", "password_hash", "display_name", "role", "created_at", "updated_at",
		})
	})

	t.Run("reviews table has correct columns", func(t *testing.T) {
		verifyColumns(t, "reviews", []string{
			"id", "tool_id", "user_id", "rating_overall", "rating_ease_of_use",
			"rating_value", "rating_accuracy", "rating_speed", "rating_support",
			"pros", "cons", "primary_use_case", "reviewer_role", "company_size",
			"usage_context", "helpful_count", "moderation_status", "moderated_by",
			"moderated_at", "created_at", "updated_at",
		})
	})

	t.Run("bookmarks table has correct columns", func(t *testing.T) {
		verifyColumns(t, "bookmarks", []string{
			"id", "user_id", "tool_id", "created_at",
		})
	})
}
