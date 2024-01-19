CREATE INDEX IF NOT EXISTS cheatcodes_code_idx ON cheatcodes USING GIN (to_tsvector('simple', code));
CREATE INDEX IF NOT EXISTS cheatcodes_description_idx ON cheatcodes USING GIN (to_tsvector('simple', description));
CREATE INDEX IF NOT EXISTS cheatcodes_tags_idx ON cheatcodes USING GIN (tags);
