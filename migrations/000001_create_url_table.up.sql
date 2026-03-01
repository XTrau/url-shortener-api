CREATE TABLE urls (
	id SERIAL PRIMARY KEY,
	url TEXT UNIQUE NOT NULL,
	slug TEXT UNIQUE NOT NULL
);

CREATE INDEX idx_urls_slug ON urls (slug);
CREATE INDEX idx_urls_url ON urls (url);