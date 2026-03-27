DROP INDEX idx_urls_slug_slug_id;
CREATE INDEX idx_urls_slug ON urls (slug);

ALTER TABLE urls DROP COLUMN slug_id;