package domain

type Url struct {
	Url string `json:"url"`
}

type Slug struct {
	Slug   string `json:"slug"`
	SlugID int    `json:"slug_id"`
}
