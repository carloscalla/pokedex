package catalog

const PageSize = 20

// Offset converts a 1-based page number into a SQL OFFSET, clamping page < 1 to 1.
func Offset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}

// TotalPages computes how many pages of pageSize cover totalCount rows.
func TotalPages(totalCount, pageSize int) int {
	if totalCount <= 0 {
		return 0
	}
	pages := totalCount / pageSize
	if totalCount%pageSize != 0 {
		pages++
	}
	return pages
}
