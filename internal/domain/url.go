package domain

type Url struct {
	Url string `json:"url"`
}

type Slug struct {
	SlugID int    `json:"slug_id"`
	Text   string `json:"slug"`
}
