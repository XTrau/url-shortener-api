package domain

type Url struct {
	Url string `json:"url"`
}

type Slug struct {
	SlugID int    `json:"slug_id"`
	Text   string `json:"slug"`
}

type UrlRepository interface {
	Create(url Url, slug Slug) error
	GetUrlBySlug(slug Slug) (Url, error)
	GetSlugByUrl(url Url) (Slug, error)
	GetFreeSlugID(slugText string) (int, error)
}
