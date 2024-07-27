CREATE TABLE IF NOT EXISTS images (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "title" VARCHAR(255) NOT NULL,
  "author" VARCHAR(50),
  "description" VARCHAR(500),
  "url" VARCHAR(512) NOT NULL,
  "likes" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);
-- create above / drop below --
DROP TABLE IF EXISTS images;