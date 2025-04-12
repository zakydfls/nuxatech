package response

type Metadata struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
}
