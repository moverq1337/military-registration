CREATE TABLE IF NOT EXISTS users
(
    id        VARCHAR(255) NOT NULL PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    surname   VARCHAR(255) NOT NULL,
    is_locked BOOLEAN      NOT NULL DEFAULT FALSE,
    verdict   BOOLEAN               DEFAULT NULL
);
