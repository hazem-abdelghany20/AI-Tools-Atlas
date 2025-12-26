### Story 3.2: Tool Alternatives Backend API

**As a** developer,
**I want** an endpoint for tool alternatives and similar tools,
**So that** users can discover related options.

**Acceptance Criteria:**

**Given** I need alternatives data for a tool
**When** I implement the alternatives endpoint
**Then** the following endpoint works:

**Endpoint: `GET /api/v1/tools/:slug/alternatives`**
- Returns list of alternative and similar tools
- Response format: `{ data: { similar: [...tools], alternatives: [...tools] } }`
- Each tool in response includes: id, slug, name, logo_url, tagline, best_for, avg_rating_overall, review_count, pricing_summary (lightweight version, same as ToolCard)
- Empty arrays if no alternatives found (not an error)

**Logic:**
- Query `tool_alternatives` table where `tool_id = current_tool.id`
- Separate results by `relationship_type`:
  - `similar`: Tools tagged as similar (e.g., same category, overlapping features)
  - `alternatives`: Direct alternatives (e.g., "Use X instead of Y")
- Limit 6 per type for initial display
- Only return non-archived tools

**And:**
- Repository function: `GetToolAlternatives(toolID uint) (similar []Tool, alternatives []Tool, error)`
- Join query to load full tool objects for alternative_tool_id references
- Service layer calls repository and formats response

**Technical Implementation:**

- GORM query example:
  ```go
  var altRecords []ToolAlternative
  db.Where("tool_id = ?", toolID).
     Preload("AlternativeTool.PrimaryCategory").
     Preload("AlternativeTool.Tags").
     Find(&altRecords)

  for _, record := range altRecords {
      if record.RelationshipType == "similar" {
          similar = append(similar, record.AlternativeTool)
      } else if record.RelationshipType == "alternative" {
          alternatives = append(alternatives, record.AlternativeTool)
      }
  }
  ```
- Filter out archived alternatives: check `AlternativeTool.ArchivedAt IS NULL`

**Prerequisites:** Story 3.1 (Tool profile API)

**Files Modified:**
- `backend/internal/tools/repository.go` (add GetToolAlternatives)
- `backend/internal/tools/service.go`
- `backend/internal/tools/handler.go` (add GET /:slug/alternatives handler)

---
