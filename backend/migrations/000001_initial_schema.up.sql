-- Categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_url VARCHAR(255),
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_slug ON categories(slug);

-- Tools table
CREATE TABLE tools (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    logo_url VARCHAR(255),
    tagline VARCHAR(500),
    description TEXT,
    best_for TEXT,
    primary_use_cases TEXT,
    pricing_summary VARCHAR(500),
    target_roles TEXT,
    platforms TEXT,
    has_free_tier BOOLEAN DEFAULT false,
    official_url VARCHAR(500),
    primary_category_id INT,
    avg_rating_overall DECIMAL(3,2) DEFAULT 0,
    review_count INT DEFAULT 0,
    bookmark_count INT DEFAULT 0,
    trending_score DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP,
    FOREIGN KEY (primary_category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_tools_slug ON tools(slug);
CREATE INDEX idx_tools_primary_category_id ON tools(primary_category_id);
CREATE INDEX idx_tools_avg_rating ON tools(avg_rating_overall);
CREATE INDEX idx_tools_bookmark_count ON tools(bookmark_count);
CREATE INDEX idx_tools_trending_score ON tools(trending_score);

-- Tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tool_tags junction table
CREATE TABLE tool_tags (
    tool_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (tool_id, tag_id),
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Media table
CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    tool_id INT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('screenshot', 'video')),
    url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
);

-- Badges table
CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tool_badges junction table
CREATE TABLE tool_badges (
    tool_id INT NOT NULL,
    badge_id INT NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tool_id, badge_id),
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE
);

-- Tool_alternatives table
CREATE TABLE tool_alternatives (
    id SERIAL PRIMARY KEY,
    tool_id INT NOT NULL,
    alternative_tool_id INT NOT NULL,
    relationship_type VARCHAR(50) NOT NULL CHECK (relationship_type IN ('similar', 'alternative')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
    FOREIGN KEY (alternative_tool_id) REFERENCES tools(id) ON DELETE CASCADE
);

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    role VARCHAR(50) NOT NULL CHECK (role IN ('user', 'admin', 'moderator')) DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Reviews table
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    tool_id INT NOT NULL,
    user_id INT NOT NULL,
    rating_overall INT NOT NULL CHECK (rating_overall >= 1 AND rating_overall <= 5),
    rating_ease_of_use INT CHECK (rating_ease_of_use >= 1 AND rating_ease_of_use <= 5),
    rating_value INT CHECK (rating_value >= 1 AND rating_value <= 5),
    rating_accuracy INT CHECK (rating_accuracy >= 1 AND rating_accuracy <= 5),
    rating_speed INT CHECK (rating_speed >= 1 AND rating_speed <= 5),
    rating_support INT CHECK (rating_support >= 1 AND rating_support <= 5),
    pros TEXT,
    cons TEXT,
    primary_use_case VARCHAR(255),
    reviewer_role VARCHAR(255),
    company_size VARCHAR(100),
    usage_context TEXT,
    helpful_count INT DEFAULT 0,
    moderation_status VARCHAR(50) NOT NULL CHECK (moderation_status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
    moderated_by INT,
    moderated_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (moderated_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_reviews_tool_id ON reviews(tool_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);

-- Bookmarks table
CREATE TABLE bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    tool_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, tool_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
);

CREATE INDEX idx_bookmarks_user_tool ON bookmarks(user_id, tool_id);
