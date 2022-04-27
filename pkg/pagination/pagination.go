package pagination

type Pagination struct {
	CurrentPage  uint `json:"current_page"`
	NextPage     uint `json:"next_page"`
	TotalItems   uint `json:"total_items"`
	TotalPages   uint `json:"total_pages"`
	ItemsPerPage uint `json:"items_per_page"`
}
