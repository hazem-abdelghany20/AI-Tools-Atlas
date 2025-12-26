// Package platform contains shared platform utilities and dependencies.
// This file ensures required dependencies are available for future use.
package platform

import (
	_ "github.com/golang-jwt/jwt/v5"           // JWT authentication (required for auth implementation)
	_ "github.com/golang-migrate/migrate/v4"   // Database migrations (required for schema management)
)
