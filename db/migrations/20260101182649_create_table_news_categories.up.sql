CREATE TABLE news_categories (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  is_active boolean DEFAULT true,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_news_categories_is_active ON news_categories(is_active);
