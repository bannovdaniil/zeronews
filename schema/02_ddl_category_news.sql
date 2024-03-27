CREATE TABLE IF NOT EXISTS news_categories
(
    id          BIGSERIAL PRIMARY KEY NOT NULL,
    news_id     BIGINT                NOT NULL,
    category_id BIGINT                NOT NULL,
    CONSTRAINT news_categories_pk UNIQUE (news_id, category_id)
);
