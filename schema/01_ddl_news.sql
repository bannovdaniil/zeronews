CREATE TABLE IF NOT EXISTS news
(
    id      BIGSERIAL PRIMARY KEY,
    title   VARCHAR(128) NOT NULL,
    content VARCHAR      NOT NULL
);
