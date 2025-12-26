package tools

// ToolFilters defines the available filters for tool queries
type ToolFilters struct {
	Category  string  // Category slug to filter by
	Price     string  // free, freemium, paid
	MinRating float64 // Minimum average rating
	Platform  string  // web, mobile, api (searches in platforms text field)
	Sort      string  // top_rated, most_bookmarked, trending, newest
}

// SortOptions defines valid sort options
const (
	SortTopRated       = "top_rated"
	SortMostBookmarked = "most_bookmarked"
	SortTrending       = "trending"
	SortNewest         = "newest"
)

// ValidateSort returns a valid sort option, defaulting to top_rated
func ValidateSort(sort string) string {
	switch sort {
	case SortTopRated, SortMostBookmarked, SortTrending, SortNewest:
		return sort
	default:
		return SortTopRated
	}
}

// ValidatePrice returns a valid price filter or empty string
func ValidatePrice(price string) string {
	switch price {
	case "free", "freemium", "paid":
		return price
	default:
		return ""
	}
}
