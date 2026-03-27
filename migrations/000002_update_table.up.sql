ALTER TABLE urls ADD slug_id SERIAL;

DROP INDEX idx_urls_slug;
CREATE INDEX idx_urls_slug_slug_id ON urls (slug, slug_id);