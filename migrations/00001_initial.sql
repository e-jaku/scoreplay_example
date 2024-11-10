-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "media" (
  id uuid DEFAULT uuid_generate_v4 (),
  name VARCHAR(255),
  file_url VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE TABLE "tag" (
  id uuid DEFAULT uuid_generate_v4 (),
  name VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE TABLE "media_tag" (
    media_id UUID REFERENCES media(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tag(id) ON DELETE CASCADE,
    PRIMARY KEY (media_id, tag_id)
);

-- Create an index on tag id since it will be used for filtering
CREATE INDEX idx_media_tag_id ON "media_tag"(tag_id);


-- +goose Down
DROP INDEX idx_media_tag_id;

DROP TABLE IF EXISTS "media_tag";
DROP TABLE IF EXISTS "tag";
DROP TABLE IF EXISTS "media";

DROP EXTENSION IF EXISTS "uuid-ossp";
