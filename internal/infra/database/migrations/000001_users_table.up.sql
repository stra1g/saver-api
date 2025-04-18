CREATE TYPE "user_roles" AS ENUM (
  'ROOT',
  'ADMIN',
  'COMMON_USER'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_deleted" boolean DEFAULT false,
  "deleted_at" timestamp DEFAULT null,
  "role" user_roles NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX users_email_unique_not_deleted ON users (email)
WHERE is_deleted = false;
