CREATE TABLE IF NOT EXISTS cheatcodes (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  code text NOT NULL,
  description text,
  tags text[] NOT NULL,
  version integer NOT NULL DEFAULT 1
);

ALTER TABLE cheatcodes ADD CONSTRAINT tags_length_check CHECK (array_length(tags, 1) BETWEEN  1 AND 5);
