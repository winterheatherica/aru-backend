CREATE TABLE awards (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  image_path text NOT NULL,
  year integer,
  order_index integer DEFAULT 0,
  is_active boolean DEFAULT true,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_awards_year ON awards(year);
CREATE INDEX idx_awards_order_index ON awards(order_index);
CREATE INDEX idx_awards_is_active ON awards(is_active);
