BEGIN;

CREATE TABLE IF NOT EXISTS "users"
(
    id      uuid,
    name    VARCHAR(32) NOT NULL,
    PRIMARY KEY ("id")
);

END;