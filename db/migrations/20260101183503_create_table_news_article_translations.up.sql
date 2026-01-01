CREATE TABLE news_article_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  article_id uuid NOT NULL REFERENCES news_articles(id) ON DELETE CASCADE,
  language language NOT NULL,
  slug text NOT NULL,
  title text NOT NULL,
  content text NOT NULL,
  meta_title text,
  meta_description text,
  meta_keywords text[],
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (article_id, language),
  UNIQUE (language, slug)
);

CREATE INDEX idx_news_article_translations_language ON news_article_translations(language);
CREATE INDEX idx_news_article_translations_article ON news_article_translations(article_id);
