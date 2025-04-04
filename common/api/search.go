package api

type SearchParameters struct {
	SortBy string   `json:"sortBy"`
	Title  string   `json:"title"`
	Topics []string `json:"topics"`
}

type ArticleResult struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type SearchResult struct {
	Aricles []ArticleResult `json:"articles"`
}
