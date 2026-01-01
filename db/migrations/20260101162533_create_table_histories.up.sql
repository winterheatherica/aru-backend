CREATE TABLE histories (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  language language NOT NULL,
  year integer,
  title text,
  description text,
  table_headers text[],
  table_rows text[][],
  is_active boolean DEFAULT true,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_histories_language ON histories(language);
CREATE INDEX idx_histories_year ON histories(year);
CREATE INDEX idx_histories_is_active ON histories(is_active);
