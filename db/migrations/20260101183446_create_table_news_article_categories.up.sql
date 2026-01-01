CREATE TABLE news_article_categories (
  article_id uuid NOT NULL REFERENCES news_articles(id) ON DELETE CASCADE,
  category_id uuid NOT NULL REFERENCES news_categories(id) ON DELETE CASCADE,
  PRIMARY KEY (article_id, category_id)
);

CREATE INDEX idx_news_article_categories_category ON news_article_categories(category_id);
