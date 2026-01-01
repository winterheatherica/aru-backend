CREATE TABLE news_category_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  category_id uuid NOT NULL REFERENCES news_categories(id) ON DELETE CASCADE,
  language language NOT NULL,
  name text NOT NULL,
  slug text NOT NULL,
  description text,
  meta_title text,
  meta_description text,
  meta_keywords text[],
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (category_id, language),
  UNIQUE (language, slug)
);

CREATE INDEX idx_news_category_translations_language ON news_category_translations(language);
CREATE INDEX idx_news_category_translations_category ON news_category_translations(category_id);
