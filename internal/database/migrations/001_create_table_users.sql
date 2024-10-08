CREATE TABLE IF NOT EXISTS users (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    "username" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "bio" VARCHAR(500),
    "profile_picture_url" VARCHAR(512),
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);